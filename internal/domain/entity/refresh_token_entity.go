package entity

import "time"

type RefreshToken struct {
	Token  string
	UserID string

	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewRefreshToken(token string, userId string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		Token:  token,
		UserID: userId,

		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}

func (t *RefreshToken) isValid() bool {
	return time.Now().Before(t.ExpiresAt)
}
