package repository

import (
	"mob/ddd-template/internal/domain/entity"
	vo "mob/ddd-template/internal/domain/valueobject"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetAll(opts *vo.PaginateOptions) (*vo.PaginatedResult[*entity.User], error)
	GetByID(id string) (*entity.User, error)
	GetByRefreshToken(refreshTokenValue string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Delete(id string) error
}
