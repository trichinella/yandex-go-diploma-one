package balance

import (
	"context"
	"diploma1/internal/app/erroring"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/ctxenv"
	"github.com/google/uuid"
	"time"
)

//go:generate easyjson -disallow_unknown_fields -all ./withdraw.go

type WithdrawOrder struct {
	OrderNumber int        `json:"order"`
	Sum         float64    `json:"sum"`
	Date        *time.Time `json:"processed_at"`
}

//easyjson:json
type WithdrawOrderList []WithdrawOrder

func WithdrawalOrders(ctx context.Context, orderRepo repo.OrderRepository) (WithdrawOrderList, error) {
	userID, ok := ctx.Value(ctxenv.ContextUserID).(uuid.UUID)
	if !ok {
		return nil, erroring.ErrIncorrectUserID
	}

	orders, err := orderRepo.WithdrawalOrdersByUser(ctx, userID)

	var withdrawalOrders []WithdrawOrder
	for _, order := range orders {
		withdrawalOrders = append(withdrawalOrders, WithdrawOrder{
			OrderNumber: order.Number,
			Sum:         order.Paid,
			Date:        order.CreatedDate,
		})
	}

	return withdrawalOrders, err
}
