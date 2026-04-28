package dto

import (
	"mob/ddd-template/internal/domain/entity"
)

type UserCreateInput struct {
	Name     string
	Email    string
	Password string
}

type UserUpdateInput struct {
	Name     *string
	Email    *string
	Password *string
}

type UserOutput struct {
	ID        string
	Name      string
	Email     string
	CreatedAt string
}

func UserEntityToOutput(e *entity.User) *UserOutput {
	return &UserOutput{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		CreatedAt: e.CreatedAt.String(),
	}
}
