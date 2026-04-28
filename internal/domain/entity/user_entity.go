package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID string

	Name         string
	Email        string
	Password string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewUser(name string, email string, password string) *User {
	return &User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		Password: password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
