package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/Akihira77/go_whatsapp/src/services"
	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/Akihira77/go_whatsapp/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type MessageType string

var (
	PEER_CHAT            MessageType = "PEER_CHAT"
	GROUP_CHAT           MessageType = "GROUP_CHAT"
	MARK_MSGS_AS_READ    MessageType = "MARK_MSGS_AS_READ"
	RECEIVE_NOTIFICATION MessageType = "RECEIVE_NOTIFICATION"
	SEND_NOTIFICATION    MessageType = "SEND_NOTIFICATION"
)

type WsMessageBody struct {
	SenderID   string     `json:"senderId"`
	ReceiverID string     `json:"receiverId,omitempty"`
	GroupID    string     `json:"groupId,omitempty"`
	Content    string     `json:"content"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
}

type WsMessage struct {
	Type MessageType    `json:"type"`
	Body *WsMessageBody `json:"body,omitempty"`
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	UserID   string
	User     *types.User
	GroupIds map[string]bool
}

type Hub struct {
	sync.RWMutex
	clients    map[string]*Client
	broadcast  chan *WsMessage
	register   chan *Client
	unregister chan *Client
	v          *utils.MyValidator
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *WsMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
		v:          utils.NewMyValidator(),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Hub) Run() {
	cleanupClient := func(client *Client) {
		h.RWMutex.Lock()
		delete(h.clients, client.UserID)
		close(client.send)
		h.RWMutex.Unlock()
	}

	for {
		select {
		case client := <-h.register:
			h.RWMutex.Lock()
			h.clients[client.UserID] = client
			h.RWMutex.Unlock()

			slog.Info("Client connected to ws pool",
				"client", client.User.Email,
			)
		case client := <-h.unregister:
			if _, ok := h.clients[client.UserID]; ok {
				slog.Info("Client disconnected from ws pool",
					"client", client.User.Email,
				)

				cleanupClient(client)
			}
		case msg := <-h.broadcast:
			b, err := json.Marshal(msg)
			if err != nil {
				slog.Error("Marshalling message error",
					"error", err,
					"sender", msg.Body.SenderID,
					"receiver", msg.Body.ReceiverID,
				)
				continue
			}

			switch msg.Type {
			case PEER_CHAT:
				if client, ok := h.clients[msg.Body.SenderID]; ok {
					slog.Info("Client sending peer to peer msg",
						"sender", client.User.Email,
					)

					select {
					case client.send <- b:
					default:
						slog.Error("Client sending peer to peer msg",
							"sender", client.User.Email,
						)
						cleanupClient(client)
					}
				}

				if client, ok := h.clients[msg.Body.ReceiverID]; ok {
					slog.Info("Client receiving peer to peer msg",
						"receiver", client.User.Email,
					)

					select {
					case client.send <- b:
					default:
						slog.Error("Client receiving peer to peer msg",
							"sender", client.User.Email,
						)
						cleanupClient(client)
					}

				}
			case SEND_NOTIFICATION:
				if client, ok := h.clients[msg.Body.SenderID]; ok {
					slog.Info("Client send notification msg",
						"receiver", client.User.Email,
					)

					select {
					case client.send <- b:
					default:
						slog.Error("Client send notification msg",
							"sender", client.User.Email,
						)
						cleanupClient(client)
					}
				}
			case GROUP_CHAT:
				for userId := range h.clients {
					client := h.clients[userId]
					if _, ok := client.GroupIds[msg.Body.GroupID]; ok {
						slog.Info("Client sending to group msg",
							"client", client.User.Email,
						)

						select {
						case client.send <- b:
						default:
							cleanupClient(client)
						}
					}
				}
			default:
			}
		}
	}
}

func (c *Client) readPump(userService *services.UserService, chatService *services.ChatService) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()

		userService.UpdateUserStatus(context.Background(), c.User, types.OFFLINE)
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Reading websocket message",
					"error", err,
				)
			}
			return
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))

		var data WsMessage
		err = json.Unmarshal(msg, &data)
		if err != nil {
			slog.Error("Unmarshalling websocket message",
				"error", err,
			)
			return
		}

		slog.Info("Receive",
			"data", data,
		)

		switch data.Type {
		case MARK_MSGS_AS_READ:
			slog.Info("Marking message as read")
			_, err := chatService.MarkMessagesAsRead(context.Background(), data.Body.SenderID, c.UserID, data.Body.GroupID)
			if err != nil {
				slog.Error("Marking messages as read",
					"error", err,
				)
				return
			}
		case SEND_NOTIFICATION:
			slog.Info("Send notification",
				"sender", data.Body.SenderID,
			)
		case PEER_CHAT:
			data.Body.SenderID = c.UserID
			slog.Info("Adding message PEER CHAT")

			m, err := chatService.AddMessage(context.Background(), &types.CreateMessage{
				Content:    data.Body.Content,
				SenderID:   data.Body.SenderID,
				ReceiverID: data.Body.ReceiverID,
				GroupID:    data.Body.GroupID,
			})
			if err != nil {
				slog.Error("Adding message",
					"error", err,
				)
				return
			}

			data.Body.CreatedAt = &m.CreatedAt
		case GROUP_CHAT:
			data.Body.SenderID = c.UserID
			slog.Info("Adding message GROUP CHAT")

			m, err := chatService.AddMessage(context.Background(), &types.CreateMessage{
				Content:    data.Body.Content,
				SenderID:   data.Body.SenderID,
				ReceiverID: data.Body.ReceiverID,
				GroupID:    data.Body.GroupID,
			})
			if err != nil {
				slog.Error("Adding message",
					"error", err,
				)
				return
			}

			data.Body.CreatedAt = &m.CreatedAt
		default:
			slog.Error("Unknown message type")
			return
		}

		c.hub.broadcast <- &data
	}
}

func (c *Client) writePump(userService *services.UserService) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()

		userService.UpdateUserStatus(context.Background(), c.User, types.OFFLINE)
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				slog.Error("Sending data not ok")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Error("Failed getting next writer",
					"error", err,
				)
				return
			}

			w.Write(msg)

		Loop:
			for {
				select {
				case msg := <-c.send:
					w.Write(newline)
					w.Write(msg)
				default:
					// No more messages, finalize the writer
					break Loop
				}
			}

			if err := w.Close(); err != nil {
				slog.Error("Writer closed",
					"error", err,
				)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("Failed to pinging websocket",
					"error", err,
				)
				return
			}
		}
	}
}

func ServeWs(c *gin.Context, hub *Hub, userService *services.UserService, chatService *services.ChatService) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Failed upgrading to websocket",
			"error", err,
		)
		return
	}

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieving your user's info")
		return
	}
	slog.Info("Your user's info",
		"user", user.Email,
	)

	groupdIds := make(map[string]bool)
	groups, err := userService.FindGroups(c, user.ID)
	if err != nil {
		slog.Error("Failed retrieving your groups",
			"error", err,
		)
		return
	}

	for _, g := range groups {
		groupdIds[g.GroupID] = true
	}
	slog.Info("Your group count",
		"group", len(groups),
	)

	user, err = userService.UpdateUserStatus(c, user, types.ONLINE)
	if err != nil {
		slog.Error("Failed update your ON/OFF-Line status",
			"error", err,
		)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte),
		UserID:   user.ID,
		User:     user,
		GroupIds: groupdIds,
	}
	client.hub.register <- client

	go client.writePump(userService)
	go client.readPump(userService, chatService)
}
