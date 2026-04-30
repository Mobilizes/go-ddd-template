package entity

import "time"

type RefreshToken struct {
	Token   string
	OwnerID string

	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewRefreshToken(token string, ownerId string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		Token:   token,
		OwnerID: ownerId,

		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}

func (t *RefreshToken) isValid() bool {
	return time.Now().Before(t.ExpiresAt)
}
