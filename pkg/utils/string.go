package utils

import "fmt"

func PendingRegistrationKey(email string) string {
	return fmt.Sprintf("pending:verifications:%s", email)
}

func AllowIPKey(email, ipaddr string) string {
	return fmt.Sprintf("pending:users:%s:allowip:%s", email, ipaddr)
}
