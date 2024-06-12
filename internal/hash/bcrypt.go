package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type bcryptHash struct {
	cost int
}

func NewBryptHash(cost ...int) HashFunction {
	bcryptCost := bcrypt.DefaultCost
	if len(cost) > 0 {
		bcryptCost = cost[0]
	}
	return &bcryptHash{cost: bcryptCost}
}

func (h *bcryptHash) Hash(salt, text string) (string, error) {
	saltedText := salt + text
	hashed, err := bcrypt.GenerateFromPassword([]byte(saltedText), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash: %w", err)
	}
	return string(hashed), nil
}

func (h *bcryptHash) Verify(salt, hashed, plain string) error {
	saltedText := salt + plain
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(saltedText))
}
