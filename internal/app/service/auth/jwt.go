package auth

import (
	"diploma1/internal/app/config"
	"diploma1/internal/app/entity"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
	"time"
)

const TokenExp = time.Hour * 3

var ErrIncorrectTokenString = errors.New("incorrect token string")

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

func GenerateTokenString(token *jwt.Token) (string, error) {
	tokenString, err := token.SignedString([]byte(config.State().JWTKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WrapTokenString(tokenString string) string {
	return `Bearer ` + tokenString
}

func UnWrapTokenString(tokenString string) (string, error) {
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return "", ErrIncorrectTokenString
	}

	return strings.TrimPrefix(tokenString, "Bearer "), nil
}

func GetClaims(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.State().JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if uuid.Nil == claims.UserID {
		return nil, fmt.Errorf("empty user ID")
	}

	return claims, err
}
