package order

import (
	"context"
	"diploma1/internal/app/entity"
	"diploma1/internal/app/erroring"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/algo/luhn"
	"diploma1/internal/app/service/ctxenv"
	"github.com/google/uuid"
	"strconv"
)

func AddOrder(ctx context.Context, repository repo.OrderRepository, input []byte, orderChannel chan<- entity.Order) error {
	if len(input) == 0 {
		return erroring.ErrEmptyRequest
	}

	orderNumber, err := strconv.Atoi(string(input))
	if err != nil {
		return erroring.ErrIncorrectNumber
	}

	if !luhn.Valid(orderNumber) {
		return &luhn.LuhnNumberError{Number: orderNumber}
	}

	order, err := repository.OrderByNumber(ctx, orderNumber)
	if err != nil {
		return err
	}

	userID, ok := ctx.Value(ctxenv.ContextUserID).(uuid.UUID)
	if !ok {
		return erroring.ErrIncorrectUserID
	}

	if order != nil {
		if order.UserID != userID {
			return &erroring.SomeoneElseOrderError{
				OwnerID:     order.UserID,
				TryUserID:   userID,
				OrderNumber: order.Number,
			}
		}

		return &erroring.NumberExistsError{
			OrderNumber: order.Number,
			UserID:      userID,
		}
	}

	order = entity.NewOrder()
	order.Number = orderNumber
	order.StatusID = repository.OrderStatusByCode(entity.NEW).ID
	order.UserID = userID

	order, err = repository.AddOrder(ctx, *order)
	if err != nil {
		return err
	}

	go func() {
		orderChannel <- *order
	}()
	return nil
}
