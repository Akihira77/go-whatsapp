package handlers

import (
	"log/slog"

	"github.com/Akihira77/go_whatsapp/src/services"
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
	slog.Info("Render home")
	views.Home().Render(c, c.Writer)
}
