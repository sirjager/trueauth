package tokens

import (
	"errors"
	"time"

	"github.com/sirjager/trueauth/pkg/utils"
)

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is invalid")

type PayloadData struct {
	SID       string `json:"id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	ClientIP  string `json:"client_ip,omitempty"`
}

type Payload struct {
	ID        string      `json:"id,omitempty"`
	IssuedAt  time.Time   `json:"iat,omitempty"`
	ExpiresAt time.Time   `json:"expires,omitempty"`
	Data      PayloadData `json:"data,omitempty"`
}

func NewPayload(p PayloadData, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		ID:        p.SID, //optionally set ID
		Data:      p,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	if payload.ID == "" {
		payload.ID = utils.UUID_XID()
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
