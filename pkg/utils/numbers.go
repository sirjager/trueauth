package utils

import (
	"strconv"
)

func GetInt64OrError(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func ToInt32(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	return int32(i), err
}
