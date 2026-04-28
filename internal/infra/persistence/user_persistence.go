package persistence

import (
	"mob/ddd-template/internal/domain/entity"
	"mob/ddd-template/internal/domain/repository"
	vo "mob/ddd-template/internal/domain/valueobject"

	"gorm.io/gorm"
)

type UserPersistence struct {
	db *gorm.DB
}

func NewUserPersistence(db *gorm.DB) repository.UserRepository {
	return &UserPersistence{db: db}
}

func (p *UserPersistence) Create(user *entity.User) error {
	return p.db.Create(user).Error
}

func (p *UserPersistence) GetAll(opts *vo.PaginateOptions) (*vo.PaginatedResult[*entity.User], error) {
	var users []*entity.User
	var total int64

	if err := p.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := opts.Page * opts.Limit

	if err := p.db.Limit(opts.Limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	return &vo.PaginatedResult[*entity.User]{
		Data:      users,
		Page:      opts.Page,
		Limit:     opts.Limit,
		TotalData: total,
		TotalPage: int((total + int64(opts.Limit) - 1) / int64(opts.Limit)),
	}, nil
}

func (p *UserPersistence) GetByID(id string) (*entity.User, error) {
	var user entity.User
	if err := p.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *UserPersistence) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := p.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *UserPersistence) Delete(id string) error {
	return p.db.Delete(&entity.User{}, "id = ?", id).Error
}
