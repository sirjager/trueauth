package hash

const (
	ErrHashesDoesNotMatch = "hashes do not match"
)

// HashFunction interface for hashing and verifying hashes
type HashFunction interface {
	Hash(salt, text string) (string, error)
	Verify(salt, hashed, plain string) error
}
