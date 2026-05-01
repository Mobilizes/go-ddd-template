package port

type TokenGenerator interface {
	GenerateAccessToken(userId string) (string, error)
	GenerateRefreshToken() (string, error)
}
