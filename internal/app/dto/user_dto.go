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
	Name  string
	Email string
}

func UserEntityToOutput(e *entity.User) *UserOutput {
	return &UserOutput{
		Name:  e.Name,
		Email: e.Email,
	}
}
