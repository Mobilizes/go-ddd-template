package entity

import (
	"time"
)

type Common struct {
	ID string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
