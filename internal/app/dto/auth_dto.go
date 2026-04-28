package dto

type AuthLoginInput struct {
	Email    string
	Password string
}

type AuthLoginOutput struct {
	Email string
	Name  string
	Token string
}
