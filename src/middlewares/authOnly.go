package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Akihira77/go_whatsapp/src/services"
	"github.com/gin-gonic/gin"
)

func AuthOnly(c *gin.Context, userService *services.UserService) {
	authHeaders := strings.Split(c.GetHeader("Authorization"), " ")
	if len(authHeaders) < 2 || authHeaders[1] == "" {
		slog.Error("User unauthorized")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User unauthorized"})
		return
	}

	user, err := userService.GetMyInfo(c, authHeaders[1])
	if err != nil {
		slog.Error("Signup",
			"error", err,
		)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed retrieving your user info"})
		return
	}

	c.Set("user", user)
	c.Next()
}
