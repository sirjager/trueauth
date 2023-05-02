package tokens

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is invalid")

type PayloadData struct {
	UserID    uuid.UUID `json:"user_id,omitempty"`
	UserEmail string    `json:"user_email,omitempty"`

	RecoveryCode     string `json:"recovery_code,omitempty"`
	DeletionCode     string `json:"deletion_code,omitempty"`
	VerificationCode string `json:"email_verification_code,omitempty"`

	AllowIP     string `json:"allow_ip,omitempty"`
	AllowIPCode string `json:"allow_ip_code,omitempty"`
}

type Payload struct {
	Id        uuid.UUID   `json:"id,omitempty"`
	IssuedAt  time.Time   `json:"iat,omitempty"`
	ExpiresAt time.Time   `json:"expires,omitempty"`
	Payload   PayloadData `json:"payload,omitempty"`
}

func NewPayload(p PayloadData, duration time.Duration) (*Payload, error) {
	token_id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		Id:        token_id,
		Payload:   p,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
