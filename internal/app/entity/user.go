package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID  `json:"uuid"`
	Login       string     `json:"login"`
	Password    string     `json:"password"`
	Balance     int64      `json:"balance"`
	CreatedDate *time.Time `json:"-"`
}

func NewUser() *User {
	id := uuid.New()
	now := time.Now().UTC()

	return &User{
		ID:          id,
		Balance:     0,
		CreatedDate: &now,
	}
}
