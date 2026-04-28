package port

type TokenGenerator interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken() (string, error)
}
