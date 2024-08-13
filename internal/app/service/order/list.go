package order

import (
	"context"
	"diploma1/internal/app/erroring"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/ctxenv"
	"github.com/google/uuid"
	"time"
)

//go:generate easyjson -disallow_unknown_fields -all ./list.go

type UserOrder struct {
	Number      int        `json:"number"`
	CreatedDate *time.Time `json:"uploaded_at"`
	Status      string     `json:"status"`
	Accrual     float64    `json:"accrual,omitempty"`
}

//easyjson:json
type UserOrderList []UserOrder

func GetUserOrderList(ctx context.Context, repository repo.OrderRepository) (UserOrderList, error) {
	userID, ok := ctx.Value(ctxenv.ContextUserID).(uuid.UUID)
	if !ok {
		return nil, erroring.ErrIncorrectUserID
	}

	orders, err := repository.OrdersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	userOrderList := make([]UserOrder, 0, len(orders))
	for _, order := range orders {
		userOrderList = append(userOrderList, UserOrder{
			Number:      order.Number,
			CreatedDate: order.CreatedDate,
			Status:      repository.OrderStatusById(order.StatusId).Code,
			Accrual:     order.Accrual,
		})
	}

	return userOrderList, nil
}
