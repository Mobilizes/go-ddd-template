package service

type PasswordHasher interface {
	HashPassword(plainPassword string) (string, error)
	ComparePassword(hash string, plainPassword string) error
}

type TokenGenerator interface {
	GenerateToken(userID string) (string, error)
}
