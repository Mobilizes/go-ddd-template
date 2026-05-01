package dto

import appDto "mob/ddd-template/internal/app/dto"

const (
	MESSAGE_SUCCESS_LOGIN      = "success login"
	MESSAGE_SUCCESS_REFRESH    = "success refresh access token"
	MESSAGE_SUCCESS_LOGOUT     = "success logout"
	MESSAGE_SUCCESS_LOGOUT_ALL = "success logout all sessions"
)

type AuthLoginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (r *AuthLoginBody) ToAppInput() *appDto.AuthLoginInput {
	return &appDto.AuthLoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

type RefreshTokenBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthLoginResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func AuthLoginOutputToResponse(out *appDto.AuthLoginOutput) *AuthLoginResponse {
	return &AuthLoginResponse{
		ID:           out.ID,
		Name:         out.Name,
		Email:        out.Email,
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}
}

type AuthMeResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthRefreshResponse struct {
	AccessToken string `json:"access_token"`
}
