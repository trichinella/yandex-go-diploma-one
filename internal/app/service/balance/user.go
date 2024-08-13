package balance

import (
	"context"
	"diploma1/internal/app/erroring"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/ctxenv"
	"github.com/google/uuid"
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

	user, err := repository.UserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &FinanceState{
		Balance: user.Balance,
		Spent:   user.Spent,
	}, nil
}
