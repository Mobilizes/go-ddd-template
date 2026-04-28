package usecase

import (
	"mob/ddd-template/internal/app/dto"
	"mob/ddd-template/internal/app/port"
	"mob/ddd-template/internal/domain/repository"

	"github.com/samber/do/v2"
)

type AuthUseCase interface {
	Login(req dto.AuthLoginInput) (dto.AuthLoginOutput, error)
}

type authUseCase struct {
	userRepository repository.UserRepository
	hasher         port.Hasher
	tokenGenerator port.TokenGenerator
}

func NewAuthUseCase(i do.Injector) AuthUseCase {
	return &authUseCase{
		userRepository: do.MustInvoke[repository.UserRepository](i),
		hasher:         do.MustInvoke[port.Hasher](i),
		tokenGenerator: do.MustInvoke[port.TokenGenerator](i),
	}
}

func (uc *authUseCase) Login(req dto.AuthLoginInput) (dto.AuthLoginOutput, error) {
	user, err := uc.userRepository.GetByEmail(req.Email)
	if err != nil {
		return dto.AuthLoginOutput{}, err
	}

	if err := uc.hasher.Compare(user.Password, req.Password); err != nil {
		return dto.AuthLoginOutput{}, err
	}

	token, err := uc.tokenGenerator.GenerateAccessToken(user.ID)
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
