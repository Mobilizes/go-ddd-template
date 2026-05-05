package persistence

import (
	"mob/ddd-template/internal/domain/entity"
	"mob/ddd-template/internal/domain/repository"

	"gorm.io/gorm"
)

type RefreshTokenPersistence struct {
	db *gorm.DB
}

func NewRefreshTokenPersistence(db *gorm.DB) repository.RefreshTokenRepository {
	return &RefreshTokenPersistence{db: db}
}

func (p *RefreshTokenPersistence) Save(token *entity.RefreshToken) error {
	return p.db.Create(token).Error
}

func (p *RefreshTokenPersistence) FindByTokenValue(tokenValue string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	if err := p.db.First(&token, "token = ?", tokenValue).Error; err != nil {
		return &entity.RefreshToken{}, err
	}

	return &token, nil
}

func (p *RefreshTokenPersistence) DeleteByTokenValue(tokenValue string) error {
	var token entity.RefreshToken
	if err := p.db.First(&token, "token = ?", tokenValue).Error; err != nil {
		return err
	}

	return p.db.Delete(&token).Error
}

func (p *RefreshTokenPersistence) DeleteAllByUserId(userId string) error {
	return p.db.Where("user_id = ?", userId).Delete(&entity.RefreshToken{}).Error
}
