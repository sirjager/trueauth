package hash

import (
	"bytes"
	"math/rand"
)

const (
	alphabets = "abcdefghijklmnopqrstuvwxyz"
	numbers   = "0123456789"
	symbols   = "_%#:>,<@!`$*()"
)

var alphabetLen, numberLen, symbolLen int

func init() {
	alphabetLen = len(alphabets)
	numberLen = len(numbers)
	symbolLen = len(symbols)
}

func RandomSalt(size ...int) string {
	saltSize := 16
	if len(size) > 0 {
		value := size[0]
		if value > 0 {
			saltSize = value
		}
	}

	var buf bytes.Buffer
	for i := 0; i < saltSize; i++ {
		switch rand.Intn(3) {
		case 0:
			buf.WriteByte(alphabets[rand.Intn(alphabetLen)])
		case 1:
			buf.WriteByte(numbers[rand.Intn(numberLen)])
		case 2:
			buf.WriteByte(symbols[rand.Intn(symbolLen)])
		}
	}

	return buf.String()
}
