package port

type PasswordHasher interface {
	HashPassword(plainPassword string) (string, error)
	ComparePassword(hash string, plainPassword string) error
}
