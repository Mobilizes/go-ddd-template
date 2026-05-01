package dto

type AuthLoginInput struct {
	Email    string
	Password string
}

type AuthLoginOutput struct {
	ID           string
	Email        string
	Name         string
	AccessToken  string
	RefreshToken string
}
