package utils

import (
	"strconv"
)

// GetInt64OrError returns int64 or error
func GetInt64OrError(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// IsNumeric returns true if string is numeric
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// ToInt32 converts string to int32
func ToInt32(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	return int32(i), err
}
