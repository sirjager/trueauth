package hash

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type argon2Hasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

type Argon2HashOpts struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

const (
	defaultTime      = 1
	defaultMemory    = 64 * 1024
	defaulThreads    = 2
	defaultKeyLength = 32
)

func NewArgon2Hasher(opts Argon2HashOpts) HashFunction {
	if opts.Time == 0 {
		opts.Time = defaultTime
	}
	if opts.Memory == 0 {
		opts.Memory = defaultMemory
	}
	if opts.Threads == 0 {
		opts.Threads = defaulThreads
	}
	if opts.KeyLen == 0 {
		opts.KeyLen = defaultKeyLength
	}
	return &argon2Hasher{
		time:    opts.Time,
		memory:  opts.Memory,
		threads: opts.Threads,
		keyLen:  opts.KeyLen,
	}
}

func (a *argon2Hasher) Hash(salt, text string) (string, error) {
	hash := argon2.IDKey(
		[]byte(text),
		[]byte(salt),
		a.time,
		a.memory,
		a.threads,
		a.keyLen,
	)
	encodedHash := fmt.Sprintf("%x", hash)
	return encodedHash, nil
}

func (a *argon2Hasher) Verify(salt, hashed, plain string) error {
	hash, err := hex.DecodeString(hashed)
	if err != nil {
		return fmt.Errorf("error decoding hash: %s", err.Error())
	}
	genHash := argon2.IDKey([]byte(plain), []byte(salt), a.time, a.memory, a.threads, a.keyLen)
	if !bytes.Equal(genHash, hash) {
		return fmt.Errorf(ErrHashesDoesNotMatch)
	}
	return nil
}
