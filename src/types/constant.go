package types

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	APPLICATION_NAME          = "Go WhatsApp"
	LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
	JWT_SIGNING_METHOD        = jwt.SigningMethodHS256
)
