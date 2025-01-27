package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Akihira77/go_whatsapp/src/services"
	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (uh *UserHandler) Signup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	var data types.Signup

	err := c.Bind(&data)
	if err != nil {
		slog.Error("Binding request payload",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request payload is invalid"})
		return
	}

	user, jwt, err := uh.userService.Signup(ctx, &data)
	if err != nil {
		slog.Error("Signup",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Signup failed"})
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", jwt))
	c.Set("user", user)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup success",
		"user":    user,
	})
}

func (uh *UserHandler) Signin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	var data types.Signin

	err := c.Bind(&data)
	if err != nil {
		slog.Error("Binding request payload",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request payload is invalid"})
		return
	}

	user, jwt, err := uh.userService.Signin(ctx, &data)
	if err != nil {
		slog.Error("Signup",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Signin failed"})
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", jwt))
	c.Set("user", user)
	c.JSON(http.StatusOK, gin.H{
		"message": "Signin success",
		"user":    user,
	})
}

func (uh *UserHandler) GetMyInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	authHeaders := strings.Split(c.GetHeader("Authorization"), " ")
	if len(authHeaders) < 2 || authHeaders[1] == "" {
		slog.Error("User unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User unauthorized"})
		return
	}

	user, err := uh.userService.GetMyInfo(ctx, authHeaders[1])
	if err != nil {
		slog.Error("Signup",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	c.Set("user", user)
	c.JSON(http.StatusOK, gin.H{
		"message": "Retrieving your user info success",
		"user":    user,
	})
}

func (uh *UserHandler) GetMyContacts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	users, err := uh.userService.GetMyContacts(ctx, user.ID)
	if err != nil {
		slog.Error("Failed retrieve user's contacts",
			"error", err,
		)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Retrieving your contacts", "users": users})
}

func (uh *UserHandler) GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	var query types.UserQuerySearch
	if err := c.ShouldBindQuery(&query); err != nil {
		slog.Error("Failed extract query",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving query search"})
		return
	}

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	users, err := uh.userService.GetUsers(ctx, user, &query)
	if err != nil {
		slog.Error("Failed retrieve user's contacts",
			"error", err,
		)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Retrieving all users", "users": users})
}
