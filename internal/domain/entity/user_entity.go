package entity

import (
	"time"
)

type User struct {
	ID string

	Name     string
	Email    string
	Password string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewUser(id string, name string, email string, password string) *User {
	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
