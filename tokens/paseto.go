package tokens

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoBuilder struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoBuilder(symmetricKey string) (TokenBuilder, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	builder := &PasetoBuilder{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return builder, nil
}

func (builder *PasetoBuilder) CreateToken(p PayloadData, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(p, duration)
	if err != nil {
		return "", payload, err
	}
	token, err := builder.paseto.Encrypt(builder.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}
	return token, payload, err
}

// Verifys token integrity and expiration time
func (builder *PasetoBuilder) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := builder.paseto.Decrypt(token, builder.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
