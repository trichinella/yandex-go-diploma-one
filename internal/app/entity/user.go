package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID  `json:"uuid"`
	Login       string     `json:"login"`
	Password    string     `json:"password"`
	Balance     float64    `json:"balance"`
	Spent       float64    `json:"spent"`
	CreatedDate *time.Time `json:"-"`
}

func NewUser() *User {
	id := uuid.New()
	now := time.Now().UTC()

	return &User{
		ID:          id,
		CreatedDate: &now,
	}
}
