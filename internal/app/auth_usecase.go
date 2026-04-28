package app

import (
	"mob/ddd-template/internal/app/dto"
	"mob/ddd-template/internal/domain/repository"
	"mob/ddd-template/internal/domain/service"

	"github.com/samber/do/v2"
)

type AuthUseCase interface {
	Login(req dto.AuthLoginInput) (dto.AuthLoginOutput, error)
}

type authUseCase struct {
	userRepository repository.UserRepository
	hasher         service.PasswordHasher
	tokenGenerator service.TokenGenerator
}

func NewAuthUseCase(i do.Injector) AuthUseCase {
	return &authUseCase{
		userRepository: do.MustInvoke[repository.UserRepository](i),
		hasher:         do.MustInvoke[service.PasswordHasher](i),
		tokenGenerator: do.MustInvoke[service.TokenGenerator](i),
	}
}

func (uc *authUseCase) Login(req dto.AuthLoginInput) (dto.AuthLoginOutput, error) {
	user, err := uc.userRepository.GetByEmail(req.Email)
	if err != nil {
		return dto.AuthLoginOutput{}, err
	}

	if err := uc.hasher.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return dto.AuthLoginOutput{}, err
	}

	token, err := uc.tokenGenerator.GenerateToken(user.ID)
	if err != nil {
		return dto.AuthLoginOutput{}, err
	}

	out := dto.AuthLoginOutput{
		Email: user.Email,
		Name:  user.Name,
		Token: token,
	}

	return out, nil
}
