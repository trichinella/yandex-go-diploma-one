package auth

import (
	"diploma1/internal/app/entity"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/secret"
	"errors"
	"fmt"
	"github.com/mailru/easyjson"
	"strings"
)

//go:generate easyjson -disallow_unknown_fields -all ./register.go

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var ErrLoginExists = errors.New("such login already exists")
var ErrBadJson = errors.New("bad json")

func Register(repository repo.UserRepositoryInterface, input []byte) error {
	registerRequest := &RegisterRequest{}
	err := easyjson.Unmarshal(input, registerRequest)

	if err != nil {
		return ErrBadJson
	}

	registerRequest.Login = strings.Trim(registerRequest.Login, " ")
	registerRequest.Password = strings.Trim(registerRequest.Password, " ")

	if registerRequest.Login == "" || registerRequest.Password == "" {
		return fmt.Errorf("login or password is empty")
	}

	user, err := repository.UserByLogin(registerRequest.Login)
	if err != nil {
		return err
	}

	if user != nil {
		return ErrLoginExists
	}

	user = entity.NewUser()
	user.Login = registerRequest.Login
	passwordHash, err := secret.HashPassword(registerRequest.Password)
	if err != nil {
		return err
	}

	user.Password = passwordHash

	user, err = repository.AddUser(*user)
	if err != nil {
		return err
	}

	return nil
}
