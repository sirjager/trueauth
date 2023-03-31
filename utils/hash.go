package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashString(text string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash : %w", err)
	}
	return string(hashed), nil
}

func VerifyHash(normal string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
}
