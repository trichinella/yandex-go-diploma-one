package balance

import (
	"context"
	"diploma1/internal/app/entity"
	"diploma1/internal/app/erroring"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/algo/luhn"
	"diploma1/internal/app/service/ctxenv"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"strconv"
)

//go:generate easyjson -disallow_unknown_fields -all ./user.go

type FinanceState struct {
	Balance float64 `json:"current"`
	Spent   float64 `json:"withdrawn"`
}

func GetUserFinanceState(ctx context.Context, repository repo.UserRepository) (*FinanceState, error) {
	userID, ok := ctx.Value(ctxenv.ContextUserID).(uuid.UUID)
	if !ok {
		return nil, erroring.ErrIncorrectUserID
	}

	user, err := repository.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &FinanceState{
		Balance: user.Balance,
		Spent:   user.Spent,
	}, nil
}

type WithdrawRequest struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}

func Withdraw(ctx context.Context, userRepo repo.UserRepository, orderRepo repo.OrderRepository, content []byte) error {
	withdrawRequest := &WithdrawRequest{}
	err := easyjson.Unmarshal(content, withdrawRequest)

	if err != nil {
		return erroring.ErrBadJSON
	}

	orderNumber, err := strconv.Atoi(withdrawRequest.OrderNumber)
	if err != nil {
		return erroring.ErrIncorrectNumber
	}

	if !luhn.Valid(orderNumber) {
		return &luhn.LuhnNumberError{Number: orderNumber}
	}

	userID, ok := ctx.Value(ctxenv.ContextUserID).(uuid.UUID)
	if !ok {
		return erroring.ErrIncorrectUserID
	}

	order, err := orderRepo.OrderByNumber(ctx, orderNumber)
	if err != nil {
		return err
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

	//mutex @todo
	user, err := userRepo.UserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Balance < withdrawRequest.Sum {
		return &erroring.LackMoneyError{
			Balance: user.Balance,
			Want:    withdrawRequest.Sum,
		}
	}

	user.Balance = user.Balance - withdrawRequest.Sum
	user.Spent = user.Spent + withdrawRequest.Sum
	order = entity.NewOrder()
	order.Number = orderNumber
	order.Paid = withdrawRequest.Sum
	order.UserID = user.ID
	order.StatusID = orderRepo.OrderStatusByCode(entity.NEW).ID

	err = userRepo.SaveUser(ctx, user)
	if err != nil {
		return err
	}

	_, err = orderRepo.AddOrder(ctx, *order)
	if err != nil {
		return err
	}

	//mutex @todo

	return nil
}
