package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/Akihira77/go_whatsapp/src/types"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(user *types.User) (string, error) {
	claims := types.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    types.APPLICATION_NAME,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(types.LOGIN_EXPIRATION_DURATION)),
		},
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.NewWithClaims(types.JWT_SIGNING_METHOD, claims)
	return token.SignedString([]byte(jwtKey))

}

func VerifyingJWT(tokenString string) (*jwt.Token, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signature")
		}

		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed parsing jwt token with claims")
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is invalid")
	}

	return token, nil
}
