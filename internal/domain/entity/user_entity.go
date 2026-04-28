package entity

import (
	"time"

	"github.com/google/uuid"
)

var AllowedUserFilters = map[string]bool{
	"name":  true,
	"email": true,
}

var AllowedUserSorts = map[string]bool{
	"name":       true,
	"created_at": true,
	"updated_at": true,
}

type User struct {
	ID string

	Name         string
	Email        string
	PasswordHash string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewUser(name string, email string, passwordHash string) *User {
	return &User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
