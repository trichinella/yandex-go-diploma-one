package repo

import (
	"context"
	"diploma1/internal/app/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	AddUser(ctx context.Context, user entity.User) (*entity.User, error)
	UserByLogin(ctx context.Context, login string) (*entity.User, error)
	UserById(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

type OrderRepository interface {
	AddOrder(ctx context.Context, order entity.Order) (*entity.Order, error)
	OrderByNumber(ctx context.Context, orderNumber int) (*entity.Order, error)
	OrderStatusByCode(statusCode entity.StatusCode) entity.OrderStatus
	OrderStatusById(statusId uuid.UUID) entity.OrderStatus
	OrdersByUser(ctx context.Context, id uuid.UUID) ([]entity.Order, error)
}
