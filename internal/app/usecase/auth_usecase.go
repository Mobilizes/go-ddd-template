package usecase

import (
	"mob/ddd-template/internal/app/dto"
	apperror "mob/ddd-template/internal/app/error"
	"mob/ddd-template/internal/app/port"
	"mob/ddd-template/internal/domain/entity"
	"mob/ddd-template/internal/domain/repository"
	"time"

	"github.com/samber/do/v2"
)

type AuthUseCase interface {
	Login(req dto.AuthLoginInput) (dto.AuthLoginOutput, error)
	Refresh(refreshToken string) (string, error)
}

type authUseCase struct {
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	hasher                 port.Hasher
	tokenGenerator         port.TokenGenerator
}

func NewAuthUseCase(i do.Injector) AuthUseCase {
	return &authUseCase{
		userRepository:         do.MustInvoke[repository.UserRepository](i),
		refreshTokenRepository: do.MustInvoke[repository.RefreshTokenRepository](i),
		hasher:                 do.MustInvoke[port.Hasher](i),
		tokenGenerator:         do.MustInvoke[port.TokenGenerator](i),
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

	refreshToken, err := uc.tokenGenerator.GenerateRefreshToken()
	if err != nil {
		return dto.AuthLoginOutput{}, err
	}

	refreshTokenEntity := entity.NewRefreshToken(refreshToken, user.ID, time.Now().Add(time.Hour*24))
	if err := uc.refreshTokenRepository.Save(refreshTokenEntity); err != nil {
		return dto.AuthLoginOutput{}, err
	}

	accessToken, err := uc.tokenGenerator.GenerateAccessToken(user.ID)
	if err != nil {
		return dto.AuthLoginOutput{}, err
	}

	out := dto.AuthLoginOutput{
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return out, nil
}

func (uc *authUseCase) Refresh(refreshToken string) (string, error) {
	_, err := uc.refreshTokenRepository.FindByTokenValue(refreshToken)
	if err != nil {
		return "", apperror.ErrRefreshTokenExpiredOrNotFound
	}

	user, err := uc.userRepository.GetByRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	accessToken, err := uc.tokenGenerator.GenerateAccessToken(user.ID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
