package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidTableName     = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]{0,62}$`).MatchString
	isValidNumber        = regexp.MustCompile(`^[0-9]+$`).MatchString
	isValidUsername      = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidName          = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString
	isValidIndentifier   = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]{0,60}$`).MatchString
	isAlphaNumUnderscore = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidEmail         = regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	).MatchString
	isValidUUID = regexp.MustCompile(
		"^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$",
	).MatchString
)

// ValidateLength validates string length
func ValidateLength(name, value string, min int, max int) error {
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
func ValidateUsername(value string) error {
	if err := ValidateLength("username", value, 3, 60); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("only be letters, digits and underscore are allowed")
	}
	return nil
}

// ValidatePassword validates password min=6 max=255
func ValidatePassword(value string) error {
	if err := ValidateLength("password", value, 8, 255); err != nil {
		return err
	}
	return nil
}

// ValidateNumber validates number
func ValidateNumber(value string) error {
	if err := ValidateLength("number", value, 1, 255); err != nil {
		return err
	}
	if !isValidNumber(value) {
		return fmt.Errorf("only be digits are allowed")
	}
	return nil
}

// ValidateEmail validates email address
func ValidateEmail(value string) error {
	if err := ValidateLength("email", value, 3, 255); err != nil {
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
func ValidateName(value string) error {
	if err := ValidateLength("name", value, 1, 255); err != nil {
		return err
	}
	if !isValidName(value) {
		return fmt.Errorf("only be letters, digits, spaces are allowed")
	}
	return nil
}

// ValidateUUID validates a UUID string
func ValidateUUID(value string) error {
	if err := ValidateLength("uuid", value, 1, 255); err != nil {
		return err
	}
	if !isValidUUID(value) {
		return fmt.Errorf("%s is not a valid uuid", value)
	}
	return nil
}

// For validating tablename min=2 max=60 letters, digits, underscore
func ValidateTableName(value string) error {
	if !isValidTableName(value) {
		return fmt.Errorf(
			`invalid schema name: %s, must be less than 62 letters, digits, underscore, starting with [a-z]`,
			value,
		)
	}
	return nil
}

// For validating tablename min=2 max=60 letters, digits, underscore
func ValidateColumnName(value string) error {
	if err := ValidateLength("column name", value, 1, 60); err != nil {
		return err
	}
	if !isAlphaNumUnderscore(value) {
		return fmt.Errorf("only be letters, digits and underscore are allowed")
	}
	return nil
}

// For validating tablename min=2 max=60 letters, digits, underscore
func ValidateIsAlphaNumUnderscore(value string) error {
	if !isAlphaNumUnderscore(value) {
		return fmt.Errorf("only be letters, digits and underscore are allowed")
	}
	return nil
}
