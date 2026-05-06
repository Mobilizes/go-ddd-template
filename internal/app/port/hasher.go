package port

type Hasher interface {
	RandomHash(plain string) (string, error)
	DeterministicHash(plain string) string
	Compare(hash string, plain string) error
}
