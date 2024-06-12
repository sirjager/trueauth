package server

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidNumber        = regexp.MustCompile(`^[0-9]+$`).MatchString
	isValidUsername      = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidName          = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString
	isAlphaNumUnderscore = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidEmail         = regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	).MatchString
)

// ValidateLength validates string length
func validateLength(name, value string, min int, max int) error {
	n := len(value)
	if n < min || n > max {
		return fmt.Errorf(
			"%s must be between %d-%d characters but got only %d chars",
			name,
			min,
			max,
			len(value),
		)
	}
	return nil
}

// ValidateUsername validates username min=3 max=60 letters, digits, underscore
func validateUsername(value string) error {
	if err := validateLength("username", value, 3, 60); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("only be letters, digits and underscore are allowed")
	}
	return nil
}

// ValidatePassword validates password min=6 max=255
func validatePassword(value string) error {
	if err := validateLength("password", value, 8, 255); err != nil {
		return err
	}
	return nil
}

// ValidateNumber validates number
func validateNumber(value string) error {
	if err := validateLength("number", value, 1, 255); err != nil {
		return err
	}
	if !isValidNumber(value) {
		return fmt.Errorf("only be digits are allowed")
	}
	return nil
}

// ValidateEmail validates email address
func validateEmail(value string) error {
	if err := validateLength("email", value, 3, 255); err != nil {
		return err
	}
	if !isValidEmail(value) {
		return fmt.Errorf("not a valid email address")
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("not a valid email address")
	}
	return nil
}

// ValidateName validates names min=1 max=255
func validateName(value string) error {
	if err := validateLength("name", value, 1, 255); err != nil {
		return err
	}
	if !isValidName(value) {
		return fmt.Errorf("only be letters, digits, spaces are allowed")
	}
	return nil
}

// For validating tablename min=2 max=60 letters, digits, underscore
func validateIsAlphaNumUnderscore(value string) error {
	if !isAlphaNumUnderscore(value) {
		return fmt.Errorf("only be letters, digits and underscore are allowed")
	}
	return nil
}
