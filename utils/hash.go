package utils

import (
	"fmt"
	"io"

	"crypto/md5"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(text string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash : %w", err)
	}
	return string(hashed), nil
}

func VerifyPassword(normal string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
}

func HashMD5(text string) (string, error) {
	h := md5.New()
	if _, err := h.Write([]byte(text)); err != nil {
		return "", err
	}
	hashed := h.Sum(nil)
	return fmt.Sprintf("%x", hashed), nil
}

func VerifyHashMD5(plainText string, hashValue string) bool {
	h := md5.New()
	io.WriteString(h, plainText)
	calculatedHash := fmt.Sprintf("%x", h.Sum(nil))
	return calculatedHash == hashValue
}
