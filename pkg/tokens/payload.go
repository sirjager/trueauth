package tokens

import (
	"errors"
	"time"

	"github.com/sirjager/trueauth/pkg/utils"
)

var (
	// ErrExpiredToken is returned when a token has expired
	ErrExpiredToken = errors.New("token expired")

	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")
)

// PayloadData contains the payload data of the token
type PayloadData struct {
	Code      string `json:"code,omitempty"`
	Type      string `json:"type,omitempty"`
	UserID    []byte `json:"user_id,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	ClientIP  string `json:"client_ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// Payload contains the payload data of the token
type Payload struct {
	IssuedAt  time.Time   `json:"iat,omitempty"`
	ExpiresAt time.Time   `json:"expires,omitempty"`
	Payload   PayloadData `json:"payload,omitempty"`
	ID        []byte      `json:"id,omitempty"`
}

// NewPayload creates a new payload for a specific username and duration
func NewPayload(p PayloadData, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		Payload:   p,
		IssuedAt:  time.Now(),
		ID:        utils.XIDNew().Bytes(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is not expired
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
