package auth

import (
	"context"
	"diploma1/internal/app/entity"
	"diploma1/internal/app/erroring"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/logging"
	"diploma1/internal/app/service/secret"
	"fmt"
	"github.com/mailru/easyjson"
	"strings"
)

//go:generate easyjson -disallow_unknown_fields -all ./register.go

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginExistsError struct {
	Login string
}

func (e *LoginExistsError) Error() string {
	return fmt.Sprintf("login already exists: \"%s\"", e.Login)
}

func Register(ctx context.Context, repository repo.UserRepository, input []byte) (string, error) {
	registerRequest := &RegisterRequest{}
	err := easyjson.Unmarshal(input, registerRequest)

	if err != nil {
		return "", erroring.ErrBadJson
	}

	registerRequest.Login = strings.Trim(registerRequest.Login, " ")
	registerRequest.Password = strings.Trim(registerRequest.Password, " ")

	if registerRequest.Login == "" || registerRequest.Password == "" {
		return "", fmt.Errorf("login or password is empty")
	}

	user, err := repository.UserByLogin(ctx, registerRequest.Login)
	if err != nil {
		return "", err
	}

	if user != nil {
		return "", &LoginExistsError{Login: user.Login}
	}

	user = entity.NewUser()
	user.Login = registerRequest.Login
	passwordHash, err := secret.HashPassword(registerRequest.Password)
	if err != nil {
		return "", err
	}

	user.Password = passwordHash

	user, err = repository.AddUser(ctx, *user)
	if err != nil {
		return "", err
	}

	token := GenerateToken(*user)
	tokenString, err := GenerateTokenString(token)
	if err != nil {
		return "", err
	}

	logging.Sugar.Infof("User with ID: %s registered successfully", user.ID)
	return tokenString, nil
}
