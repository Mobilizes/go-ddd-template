package entity

type User struct {
	Common

	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

