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
	Login(req *dto.AuthLoginInput) (*dto.AuthLoginOutput, error)
	Refresh(refreshToken string) (string, error)
	Logout(refreshToken string) error
	LogoutAll(userId string) error
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

func (uc *authUseCase) Login(req *dto.AuthLoginInput) (*dto.AuthLoginOutput, error) {
	user, err := uc.userRepository.GetByEmail(req.Email)
	if err != nil {
		return &dto.AuthLoginOutput{}, apperror.ErrInvalidEmailOrPassword
	}

	if err := uc.hasher.Compare(user.Password, req.Password); err != nil {
		return &dto.AuthLoginOutput{}, apperror.ErrInvalidEmailOrPassword
	}

	refreshToken, err := uc.tokenGenerator.GenerateRefreshToken()
	if err != nil {
		return &dto.AuthLoginOutput{}, err
	}

	hashedRefreshToken, err := uc.hasher.Hash(refreshToken)
	if err != nil {
		return &dto.AuthLoginOutput{}, err
	}

	refreshTokenEntity := entity.NewRefreshToken(hashedRefreshToken, user.ID, time.Now().Add(time.Hour*24))
	if err := uc.refreshTokenRepository.Save(refreshTokenEntity); err != nil {
		return &dto.AuthLoginOutput{}, err
	}

	accessToken, err := uc.tokenGenerator.GenerateAccessToken(user.ID)
	if err != nil {
		return &dto.AuthLoginOutput{}, err
	}

	out := dto.AuthLoginOutput{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &out, nil
}

func (uc *authUseCase) Refresh(refreshToken string) (string, error) {
	hashedRefreshToken, err := uc.hasher.Hash(refreshToken)
	if err != nil {
		return "", err
	}

	token, err := uc.refreshTokenRepository.FindByTokenValue(hashedRefreshToken)
	if err != nil {
		return "", apperror.ErrRefreshTokenExpiredOrNotFound
	}

	user, err := uc.userRepository.GetById(token.OwnerID)
	if err != nil {
		return "", err
	}

	accessToken, err := uc.tokenGenerator.GenerateAccessToken(user.ID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (uc *authUseCase) Logout(refreshToken string) error {
	hashedRefreshToken, err := uc.hasher.Hash(refreshToken)
	if err != nil {
		return err
	}

	if err := uc.refreshTokenRepository.DeleteByTokenValue(hashedRefreshToken); err != nil {
		return err
	}

	return nil
}

func (uc *authUseCase) LogoutAll(userId string) error {
	if err := uc.refreshTokenRepository.DeleteAllByUserId(userId); err != nil {
		return err
	}

	return nil
}
