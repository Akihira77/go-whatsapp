package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Akihira77/go_whatsapp/src/services"
	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/Akihira77/go_whatsapp/src/views"
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
	ctx, cancel := context.WithTimeout(c, 2*time.Second)
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

	file, _, err := c.Request.FormFile("image")
	if err != nil && file != nil {
		slog.Error("Failed extract image payload")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request payload is invalid"})
		return
	}

	image := file
	user, jwt, err := uh.userService.Signup(ctx, &data, image)
	if err != nil {
		slog.Error("Signup",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Signup failed"})
		return
	}

	c.Set("user", user)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", jwt, 60*60*24, "/", "localhost", true, true)
	c.Header("HX-Redirect", "/")

	views.Home().Render(c, c.Writer)
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
		slog.Error("Signin",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Signin failed"})
		return
	}

	c.Set("user", user)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", jwt, 60*60*24, "/", "localhost", true, true)
	c.Header("HX-Redirect", "/")

	views.Home().Render(c, c.Writer)
}

func (uh *UserHandler) GetMyImageProfile(c *gin.Context) {
	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Writer.Write(user.ImageUrl)
}

func (uh *UserHandler) GetMyInfo(c *gin.Context) {
	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

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

	search := c.Query("search")
	users, err := uh.userService.GetMyContacts(ctx, user.ID, search)
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

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	var query types.UserQuerySearch
	if err := c.ShouldBindQuery(&query); err != nil {
		slog.Error("Failed extract query",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving query search"})
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

func (uh *UserHandler) UpdatePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	var data types.UpdatePassword
	err := c.ShouldBind(&data)
	if err != nil {
		slog.Error("Failed extract request payload")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request payload is invalid"})
		return
	}

	user, err = uh.userService.UpdatePassword(ctx, user, data)
	if err != nil {
		slog.Error("Failed update user's password",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed updating your password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update profile success",
		"user":    user,
	})
}

func (uh *UserHandler) UpdateUserProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	file, _, err := c.Request.FormFile("image")
	if err != nil && file != nil {
		slog.Error("Failed extract image payload")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request payload is invalid"})
		return
	}

	var data types.UpdateUser
	err = c.ShouldBind(&data)
	if err != nil {
		slog.Error("Failed extract request payload")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request payload is invalid"})
		return
	}

	image := file
	user, err = uh.userService.UpdateUserProfile(ctx, user, &data, image)
	if err != nil {
		slog.Error("Failed update user's profile",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed updating your profile"})
		return
	}

	c.Set("user", user)
	c.Header("HX-Redirect", "/users/profile")

	views.MyProfile().Render(c, c.Writer)
}

func (uh *UserHandler) AddContact(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	userID := c.Param("userId")
	users, err := uh.userService.AddContact(ctx, user, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Adding contact failed",
			"users":   users,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Adding contact success",
		"users":   users,
	})
}

func (uh *UserHandler) RemoveContact(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	userID := c.Param("userId")
	users, err := uh.userService.RemoveContact(ctx, user, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Removing contact failed",
			"users":   users,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Removing contact success",
		"users":   users,
	})
}
