package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Akihira77/go_whatsapp/src/components"
	"github.com/Akihira77/go_whatsapp/src/services"
	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *services.ChatService
	hub         *Hub
}

func NewChatHandler(chatService *services.ChatService, hub *Hub) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		hub:         hub,
	}
}

func (ch *ChatHandler) SearchLastMessage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Retrieving user's info")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving user's info"})
		return
	}

	userName := c.Param("username")
	groupName := c.Param("groupname")

	if userName != "" {
		userName = "%" + userName + "%"
	}

	if groupName != "" {
		groupName = "%" + groupName + "%"
	}

	msgs, err := ch.chatService.SearchChat(ctx, user.ID, userName, groupName)
	if err != nil {
		slog.Error("Retrieving last messages",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving last messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Success retrieving last messages",
		"messages": msgs,
	})
}

func (ch *ChatHandler) GetChatList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Retrieving user's info")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving user's info"})
		return
	}

	userName := c.Query("username")
	groupName := c.Query("groupname")

	if userName != "" {
		userName = "%" + userName + "%"
	}

	if groupName != "" {
		groupName = "%" + groupName + "%"
	}

	users, err := ch.chatService.SearchChat(ctx, user.ID, userName, groupName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving last messages"})
		return
	}

	components.ChatList(users).Render(c, c.Writer)
}

func (ch *ChatHandler) SendMsg(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Retrieving user's info")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving user's info"})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["files[]"]

	var data types.CreateMessage
	err := c.ShouldBind(&data)
	if err != nil {
		slog.Error("Failed parsing message data from",
			"user", user.Email,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed parsing message data"})
		return
	}

	data.SenderID = user.ID
	data.Files = files
	msg, err := ch.chatService.AddMessage(ctx, &data)
	if err != nil {
		slog.Error("Failed saving chat",
			"user", user.Email,
			"err", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	go func() {
		wsMsg := &WsMessage{
			Body: &WsMessageBody{
				SenderID:   msg.SenderID,
				ReceiverID: msg.ReceiverID,
				GroupID:    msg.GroupID,
				Content:    msg.Content,
				Files:      msg.Files,
				CreatedAt:  &msg.CreatedAt,
			},
		}

		if data.GroupID != nil {
			wsMsg.Type = GROUP_CHAT
		} else {
			wsMsg.Type = PEER_CHAT
		}

		ch.hub.Broadcast <- wsMsg
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Success adding message",
		"msg":     msg,
	})
}

func (uh *ChatHandler) FindFileInsideChat(c *gin.Context) {
	msgId := c.Param("messageId")
	fileId := c.Param("fileId")
	file, err := uh.chatService.FindFileInsideChat(c, msgId, fileId)
	if err != nil {
		slog.Error("Failed retrieve chat message")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieve chat data"})
		return
	}

	mimeType := http.DetectContentType(file.Data)
	c.Header("Content-Type", mimeType)
	c.Status(http.StatusOK)
	c.Writer.Write(file.Data)
}
