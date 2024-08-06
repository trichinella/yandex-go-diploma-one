package repo

import (
	"context"
	"diploma1/internal/app/entity"
)

type UserRepositoryInterface interface {
	AddUser(ctx context.Context, user entity.User) (*entity.User, error)
	UserByLogin(ctx context.Context, login string) (*entity.User, error)
}
