package auth

import (
	"context"
	"diploma1/internal/app/erroring"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/logging"
	"diploma1/internal/app/service/secret"
	"fmt"
	"github.com/mailru/easyjson"
	"strings"
)

//go:generate easyjson -disallow_unknown_fields -all ./login.go

// LoginRequest Отдельная структура для логина и регистрации - со временем регистрация явно станет сложнее
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginAbsentError struct {
	Login string
}

func (e *LoginAbsentError) Error() string {
	return fmt.Sprintf("login absent: \"%s\"", e.Login)
}

type IncorrectPasswordError struct {
	Login    string
	Password string
}

func (e *IncorrectPasswordError) Error() string {
	return fmt.Sprintf("For login \"%s\" incorrect password: \"%s\"", e.Login, e.Password)
}

func Login(ctx context.Context, repository repo.UserRepository, input []byte) (string, error) {
	loginRequest := &LoginRequest{}
	err := easyjson.Unmarshal(input, loginRequest)

	if err != nil {
		return "", erroring.ErrBadJson
	}

	loginRequest.Login = strings.Trim(loginRequest.Login, " ")
	loginRequest.Password = strings.Trim(loginRequest.Password, " ")

	if loginRequest.Login == "" || loginRequest.Password == "" {
		return "", fmt.Errorf("login or password is empty")
	}

	user, err := repository.UserByLogin(ctx, loginRequest.Login)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", &LoginAbsentError{Login: loginRequest.Login}
	}

	if !secret.CheckPasswordHash(loginRequest.Password, user.Password) {
		return "", &IncorrectPasswordError{Login: loginRequest.Login, Password: loginRequest.Password}
	}

	token := GenerateToken(*user)
	tokenString, err := GenerateTokenString(token)
	if err != nil {
		return "", err
	}

	logging.Sugar.Infof("User with ID: %s login successfully", user.ID)

	return tokenString, nil
}
