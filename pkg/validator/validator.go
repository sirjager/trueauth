package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidName     = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString
	isValidUUID     = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$").MatchString
)

// For validating string length
func ValidateLength(name, value string, min int, max int) error {
	n := len(value)
	if n < min || n > max {
		return fmt.Errorf("%s must be between %d-%d characters but got only %d chars", name, min, max, len(value))
	}
	return nil
}

// For validating username min=3 max=60 letters, digits, underscore
func ValidateUsername(value string) error {
	if err := ValidateLength("username", value, 3, 60); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("only be letters, digits and underscore are allowed")
	}
	return nil
}

// For validating password min=6 ma=255
func ValidatePassword(value string) error {
	if err := ValidateLength("password", value, 8, 255); err != nil {
		return err
	}
	return nil
}

// For validating email address
func ValidateEmail(value string) error {
	if err := ValidateLength("email", value, 3, 255); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("not a valid email address")
	}
	return nil
}

// For validating names min=1 max=255
func ValidateName(value string) error {
	if err := ValidateLength("name", value, 1, 255); err != nil {
		return err
	}
	if !isValidName(value) {
		return fmt.Errorf("only be letters, digits, spaces are allowed")
	}
	return nil
}

// For validating uuid
func ValidateUUID(value string) error {
	if err := ValidateLength("uuid", value, 1, 255); err != nil {
		return err
	}
	if !isValidUUID(value) {
		return fmt.Errorf("%s is not a valid uuid", value)
	}
	return nil
}
