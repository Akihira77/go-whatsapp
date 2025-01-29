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

type PageHandler struct {
	userService *services.UserService
}

func NewPageHandler(userService *services.UserService) *PageHandler {
	return &PageHandler{
		userService: userService,
	}
}

func (ph *PageHandler) RenderSignup(c *gin.Context) {
	views.Signup().Render(c, c.Writer)
}

func (ph *PageHandler) RenderSignin(c *gin.Context) {
	views.Signin().Render(c, c.Writer)
}

func (ph *PageHandler) RenderHome(c *gin.Context) {
	views.Home().Render(c, c.Writer)
}

func (ph *PageHandler) RenderMyProfile(c *gin.Context) {
	c.Header("HX-Redirect", "/users/profile")

	views.MyProfile().Render(c, c.Writer)
}

func (ph *PageHandler) RenderEditProfile(c *gin.Context) {
	views.EditUser().Render(c, c.Writer)
}

func (ph *PageHandler) RenderChangePassword(c *gin.Context) {
	views.ChangePassword().Render(c, c.Writer)
}

func (ph *PageHandler) RenderMyContacts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
	defer cancel()

	user, ok := c.MustGet("user").(*types.User)
	if !ok {
		slog.Error("Failed retrieve user's data from context")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	search := c.Query("search")
	users, err := ph.userService.GetMyContacts(ctx, user.ID, search)
	if err != nil {
		slog.Error("Failed retrieve user's contacts",
			"error", err,
		)
	}

	// c.Header("HX-Redirect", "/contacts")

	views.MyContacts(user, users).Render(c, c.Writer)
}
