package repo

import (
	"context"
	"diploma1/internal/app/entity"
)

type UserRepository interface {
	AddUser(ctx context.Context, user entity.User) (*entity.User, error)
	UserByLogin(ctx context.Context, login string) (*entity.User, error)
}

type OrderRepository interface {
	AddOrder(ctx context.Context, order entity.Order) (*entity.Order, error)
	OrderByNumber(ctx context.Context, orderNumber int) (*entity.Order, error)
	OrderStatus(statusCode entity.StatusCode) entity.OrderStatus
}
