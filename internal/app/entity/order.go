package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Order struct {
	ID          uuid.UUID           `json:"uuid"`
	UserId      uuid.UUID           `json:"user_id"`
	Number      int                 `json:"number"`
	CreatedDate *time.Time          `json:"-"`
	StatusId    uuid.UUID           `json:"status_id"`
	Accrual     decimal.NullDecimal `json:"accrual"`
	Paid        decimal.Decimal     `json:"paid"`
}

func NewOrder() *Order {
	id := uuid.New()
	now := time.Now().UTC()

	return &Order{
		ID:          id,
		CreatedDate: &now,
		Paid:        decimal.Zero,
	}
}

type OrderStatus struct {
	ID    uuid.UUID `json:"uuid"`
	Title string    `json:"title"`
	Code  string    `json:"code"`
}

type StatusCode string

const (
	NEW        StatusCode = "NEW"
	PROCESSING StatusCode = "PROCESSING"
	INVALID    StatusCode = "INVALID"
	PROCESSED  StatusCode = "PROCESSED"
)
