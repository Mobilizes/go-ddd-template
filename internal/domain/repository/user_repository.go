package repository

import (
	"mob/ddd-template/internal/domain/entity"
	vo "mob/ddd-template/internal/domain/valueobject"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetAll(opts *vo.PaginateOptions) (*vo.PaginatedResult[*entity.User], error)
	GetById(id string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Delete(id string) error
}
