package service

import (
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	ClientIP  string    `json:"clientip"`
	UserAgent string    `json:"useragent"`
	UserID    string    `json:"user_id"`
	Blocked   bool      `json:"blocked"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
