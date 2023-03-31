// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package sqlc

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	// user id
	ID uuid.UUID `json:"id"`
	// unique email
	Email string `json:"email"`
	// unique username
	Username string `json:"username"`
	// hashed password
	Password string `json:"password"`
	// first name can be empty
	Firstname string `json:"firstname"`
	// last name can be empty
	Lastname string `json:"lastname"`
	// email verified or not
	Verified bool `json:"verified"`
	// user blocked or not
	Blocked bool `json:"blocked"`
	// created at timestamp
	CreatedAt time.Time `json:"created_at"`
	// last updated at timestamp
	UpdatedAt time.Time `json:"updated_at"`
}
