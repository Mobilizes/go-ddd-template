package persistence

import (
	"mob/ddd-template/internal/app/port"
	"mob/ddd-template/internal/domain/repository"

	"gorm.io/gorm"
)

type GormUnitOfWork struct {
	db *gorm.DB
}

func NewGormUnitOfWork(db *gorm.DB) port.UnitOfWork {
	return &GormUnitOfWork{db: db}
}

func (uow *GormUnitOfWork) Transaction(fn func(repos port.UnitOfWorkRepositories) error) error {
	return uow.db.Transaction(func(tx *gorm.DB) error {
		repos := &gormUnitOfWorkRepositories{db: tx}
		return fn(repos)
	})
}

type gormUnitOfWorkRepositories struct {
	db *gorm.DB
}

func (repos *gormUnitOfWorkRepositories) Users() repository.UserRepository {
	return NewUserPersistence(repos.db)
}

func (repos *gormUnitOfWorkRepositories) RefreshTokens() repository.RefreshTokenRepository {
	return NewRefreshTokenPersistence(repos.db)
}
