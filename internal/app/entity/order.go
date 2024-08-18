package entity

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID          uuid.UUID  `json:"uuid"`
	UserID      uuid.UUID  `json:"user_id"`
	Number      int        `json:"number"`
	CreatedDate *time.Time `json:"-"`
	StatusID    uuid.UUID  `json:"status_id"`
	Accrual     float64    `json:"accrual"`
	Paid        float64    `json:"paid"`
}

func NewOrder() *Order {
	id := uuid.New()
	now := time.Now().UTC()

	return &Order{
		ID:          id,
		CreatedDate: &now,
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
