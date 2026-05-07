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
	Me(userId string) (*dto.UserOutput, error)
	Refresh(refreshToken string) (string, string, error)
	Logout(userId string) error
}

type authUseCase struct {
	userRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	hasher                 port.Hasher
	tokenGenerator         port.TokenGenerator
	uow                    port.UnitOfWork
}

func NewAuthUseCase(i do.Injector) AuthUseCase {
	return &authUseCase{
		userRepository:         do.MustInvoke[repository.UserRepository](i),
		refreshTokenRepository: do.MustInvoke[repository.RefreshTokenRepository](i),
		hasher:                 do.MustInvoke[port.Hasher](i),
		tokenGenerator:         do.MustInvoke[port.TokenGenerator](i),
		uow:                    do.MustInvoke[port.UnitOfWork](i),
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

	hashedRefreshToken := uc.hasher.DeterministicHash(refreshToken)

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

func (uc *authUseCase) Me(userId string) (*dto.UserOutput, error) {
	user, err := uc.userRepository.GetById(userId)
	if err != nil {
		return &dto.UserOutput{}, apperror.ErrUserNotFound
	}

	return dto.UserEntityToOutput(user), nil
}

func (uc *authUseCase) Refresh(refreshToken string) (string, string, error) {
	var newRefreshToken string
	hashedRefreshToken := uc.hasher.DeterministicHash(refreshToken)

	token, err := uc.refreshTokenRepository.FindByTokenValue(hashedRefreshToken)
	if err != nil {
		return "", "", apperror.ErrRefreshTokenExpiredOrNotFound
	}

	user, err := uc.userRepository.GetById(token.UserID)
	if err != nil {
		return "", "", err
	}

	if err := uc.uow.Transaction(func(repos port.UnitOfWorkRepositories) error {
		refreshTokenRepo := repos.RefreshTokens()

		err = refreshTokenRepo.DeleteAllByUserId(user.ID)
		if err != nil {
			return err
		}

		newRefreshToken, err = uc.tokenGenerator.GenerateRefreshToken()
		if err != nil {
			return err
		}

		hashedRefreshToken := uc.hasher.DeterministicHash(newRefreshToken)
		refreshTokenEntity := entity.NewRefreshToken(hashedRefreshToken, user.ID, time.Now().Add(time.Hour*24))
		if err := refreshTokenRepo.Save(refreshTokenEntity); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", "", err
	}

	accessToken, err := uc.tokenGenerator.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (uc *authUseCase) Logout(userId string) error {
	return uc.uow.Transaction(func(repos port.UnitOfWorkRepositories) error {
		return repos.RefreshTokens().DeleteAllByUserId(userId)
	})
}
