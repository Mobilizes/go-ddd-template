package repository

import (
	"mob/ddd-template/internal/domain/entity"
	vo "mob/ddd-template/internal/domain/valueobject"
)

type FileRepository interface {
	Create(file *entity.File) error
	GetById(id string) (*entity.File, error)
	GetAllByUser(userId string, opts *vo.PaginateOptions) (*vo.PaginatedResult[*entity.File], error)
	Update(file *entity.File) error
	Delete(id string) error
}
