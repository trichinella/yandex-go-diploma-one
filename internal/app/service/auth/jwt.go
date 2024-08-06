package auth

import (
	"diploma1/internal/app/config"
	"diploma1/internal/app/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const TokenExp = time.Hour * 3

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

func GenerateToken(user entity.User) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: user.ID,
	})
}

func TokenRaw(token *jwt.Token) (string, error) {
	tokenString, err := token.SignedString([]byte(config.State().JWTKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
