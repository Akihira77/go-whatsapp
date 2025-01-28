package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/Akihira77/go_whatsapp/src/services"
	"github.com/gin-gonic/gin"
)

func AuthOnly(c *gin.Context, userService *services.UserService) {
	token, err := c.Cookie("token")
	if err != nil || token == "" {
		slog.Error("User unauthorized",
			"error", err,
			"token", token,
		)
		c.Redirect(http.StatusFound, "/signin")
		c.Abort()
		return
	}

	user, err := userService.GetMyInfo(c, token)
	if err != nil {
		slog.Error("Signup",
			"error", err,
		)
		c.Redirect(http.StatusFound, "/signin")
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}
