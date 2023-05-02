package crypt

import (
	"crypto/md5"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = bcrypt.DefaultCost

var md5Hasher = md5.New()

// HashPassword hashes the given text using bcrypt algorithm.
// It returns the hashed password as a string.
func HashPassword(text string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash: %w", err)
	}
	return string(hashed), nil
}

// VerifyPassword compares the normal (plain) password with the hashed password.
// It returns an error if the passwords don't match.
func VerifyPassword(normal string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
}

// HashMD5 hashes the given text using the MD5 algorithm.
// It returns the hashed text as a string.
func HashMD5(text string) (string, error) {
	_, err := io.WriteString(md5Hasher, text)
	if err != nil {
		return "", fmt.Errorf("failed to hash: %w", err)
	}
	hashed := md5Hasher.Sum(nil)
	md5Hasher.Reset()
	return fmt.Sprintf("%x", hashed), nil
}

// VerifyHashMD5 compares the plain text with the hashed value.
// It returns true if the hash values match, otherwise false.
func VerifyHashMD5(plainText string, hashValue string) bool {
	_, err := io.WriteString(md5Hasher, plainText)
	if err != nil {
		return false
	}
	calculatedHash := md5Hasher.Sum(nil)
	md5Hasher.Reset()
	return fmt.Sprintf("%x", calculatedHash) == hashValue
}
