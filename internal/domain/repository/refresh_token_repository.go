package repository

import (
	"mob/ddd-template/internal/domain/entity"
)

type RefreshTokenRepository interface {
	Save(token *entity.RefreshToken) error
	FindByTokenValue(tokenValue string) (*entity.RefreshToken, error)
	Revoke(tokenValue string) error
	DeleteAllByUserId(userId string) error
}
