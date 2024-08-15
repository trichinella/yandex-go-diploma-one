package erroring

import (
	"fmt"
	"github.com/google/uuid"
)

type SomeoneElseOrderError struct {
	OwnerID     uuid.UUID
	TryUserID   uuid.UUID
	OrderNumber int
}

func (e *SomeoneElseOrderError) Error() string {
	return fmt.Sprintf("User %s tried adding order %d (owner %s)", e.TryUserID, e.OrderNumber, e.OwnerID)
}

type NumberExistsError struct {
	OrderNumber int
	UserID      uuid.UUID
}

func (e *NumberExistsError) Error() string {
	return fmt.Sprintf("User %s has already %d)", e.UserID, e.OrderNumber)
}

type LackMoneyError struct {
	Balance float64
	Want    float64
}

func (e *LackMoneyError) Error() string {
	return fmt.Sprintf("Balance %.2f, but want to withdraw %.2f)", e.Balance, e.Want)
}
