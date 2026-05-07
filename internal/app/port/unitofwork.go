package port

import "mob/ddd-template/internal/domain/repository"

type UnitOfWork interface {
	Transaction(fn func(repos UnitOfWorkRepositories) error) error
}

type UnitOfWorkRepositories interface {
	Users() repository.UserRepository
	RefreshTokens() repository.RefreshTokenRepository
}
