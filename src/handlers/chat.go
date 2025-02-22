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
}

func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
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

func (ch *ChatHandler) SendMsgToOffUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Retrieving user's info")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving user's info"})
		return
	}

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
	msg, err := ch.chatService.AddMessage(ctx, &data)
	if err != nil {
		slog.Error("Failed saving chat",
			"user", user.Email,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed saving your chat"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Success adding message",
		"msg":     msg,
	})
}
