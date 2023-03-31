package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"
const numbers = "0123456789"

const symbols = "_%#:>,<@!`$*()"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate a random number between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Generate a random string of length len
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)
	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomNumberAsString(digits int) string {
	var sb strings.Builder
	k := len(numbers)
	for i := 0; i < digits; i++ {
		c := numbers[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// Generate a random symbols of length len
func RandomSymbols(n int) string {
	var sb strings.Builder
	k := len(symbols)
	for i := 0; i < n; i++ {
		c := symbols[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomEmail() string {
	return RandomString(5) + strconv.Itoa(int(RandomInt(1, 20))) + "@gmail.com"
}

func RandomUserName() string {
	no := int(RandomInt(5, 30))
	return RandomString(no) + strconv.Itoa(int(RandomInt(1, 20)))
}

func RandomPassword() string {
	no := int(RandomInt(8, 30))
	sym := RandomSymbols(5)
	return RandomString(no) + sym + strconv.Itoa(int(RandomInt(1, 20)))
}

func RandomTableName() string {
	no := int(RandomInt(2, 10))
	return RandomString(no) + strconv.Itoa(int(RandomInt(1, 10)))
}
