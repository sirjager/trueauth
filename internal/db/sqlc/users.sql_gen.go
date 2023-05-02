// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: users.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    email, username, password,firstname, lastname,allowed_ips,
    last_verify_sent_at,last_recovery_sent_at,last_emailchange_sent_at,last_delete_sent_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at
`

type CreateUserParams struct {
	Email                 string    `json:"email"`
	Username              string    `json:"username"`
	Password              string    `json:"password"`
	Firstname             string    `json:"firstname"`
	Lastname              string    `json:"lastname"`
	AllowedIps            []string  `json:"allowed_ips"`
	LastVerifySentAt      time.Time `json:"last_verify_sent_at"`
	LastRecoverySentAt    time.Time `json:"last_recovery_sent_at"`
	LastEmailchangeSentAt time.Time `json:"last_emailchange_sent_at"`
	LastDeleteSentAt      time.Time `json:"last_delete_sent_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.Username,
		arg.Password,
		arg.Firstname,
		arg.Lastname,
		pq.Array(arg.AllowedIps),
		arg.LastVerifySentAt,
		arg.LastRecoverySentAt,
		arg.LastEmailchangeSentAt,
		arg.LastDeleteSentAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Firstname,
		&i.Lastname,
		&i.EmailVerified,
		&i.VerifyToken,
		&i.LastVerifySentAt,
		&i.RecoveryToken,
		&i.LastRecoverySentAt,
		&i.EmailchangeToken,
		&i.LastEmailchangeSentAt,
		pq.Array(&i.AllowedIps),
		&i.AllowipToken,
		&i.DeleteToken,
		&i.LastDeleteSentAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const readUserByEmail = `-- name: ReadUserByEmail :one
SELECT id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at FROM users WHERE email = $1 LIMIT 1
`

func (q *Queries) ReadUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, readUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Firstname,
		&i.Lastname,
		&i.EmailVerified,
		&i.VerifyToken,
		&i.LastVerifySentAt,
		&i.RecoveryToken,
		&i.LastRecoverySentAt,
		&i.EmailchangeToken,
		&i.LastEmailchangeSentAt,
		pq.Array(&i.AllowedIps),
		&i.AllowipToken,
		&i.DeleteToken,
		&i.LastDeleteSentAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const readUserByID = `-- name: ReadUserByID :one
SELECT id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at FROM users WHERE id = $1 LIMIT 1
`

func (q *Queries) ReadUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, readUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Firstname,
		&i.Lastname,
		&i.EmailVerified,
		&i.VerifyToken,
		&i.LastVerifySentAt,
		&i.RecoveryToken,
		&i.LastRecoverySentAt,
		&i.EmailchangeToken,
		&i.LastEmailchangeSentAt,
		pq.Array(&i.AllowedIps),
		&i.AllowipToken,
		&i.DeleteToken,
		&i.LastDeleteSentAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const readUserByUsername = `-- name: ReadUserByUsername :one
SELECT id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at FROM users WHERE username = $1 LIMIT 1
`

func (q *Queries) ReadUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, readUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Firstname,
		&i.Lastname,
		&i.EmailVerified,
		&i.VerifyToken,
		&i.LastVerifySentAt,
		&i.RecoveryToken,
		&i.LastRecoverySentAt,
		&i.EmailchangeToken,
		&i.LastEmailchangeSentAt,
		pq.Array(&i.AllowedIps),
		&i.AllowipToken,
		&i.DeleteToken,
		&i.LastDeleteSentAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const readUsers = `-- name: ReadUsers :many
SELECT id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at FROM users LIMIT $2 OFFSET $1
`

type ReadUsersParams struct {
	Offset sql.NullInt32 `json:"offset"`
	Limit  sql.NullInt32 `json:"limit"`
}

func (q *Queries) ReadUsers(ctx context.Context, arg ReadUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, readUsers, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Username,
			&i.Password,
			&i.Firstname,
			&i.Lastname,
			&i.EmailVerified,
			&i.VerifyToken,
			&i.LastVerifySentAt,
			&i.RecoveryToken,
			&i.LastRecoverySentAt,
			&i.EmailchangeToken,
			&i.LastEmailchangeSentAt,
			pq.Array(&i.AllowedIps),
			&i.AllowipToken,
			&i.DeleteToken,
			&i.LastDeleteSentAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET
    username = $1,
    password = $2,
    firstname = $3,
    lastname = $4
    WHERE id = $5
RETURNING id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at
`

type UpdateUserParams struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Username,
		arg.Password,
		arg.Firstname,
		arg.Lastname,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Firstname,
		&i.Lastname,
		&i.EmailVerified,
		&i.VerifyToken,
		&i.LastVerifySentAt,
		&i.RecoveryToken,
		&i.LastRecoverySentAt,
		&i.EmailchangeToken,
		&i.LastEmailchangeSentAt,
		pq.Array(&i.AllowedIps),
		&i.AllowipToken,
		&i.DeleteToken,
		&i.LastDeleteSentAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserAllowIP = `-- name: UpdateUserAllowIP :exec
UPDATE users SET 
    allowed_ips = $1,
    allowip_token = $2
WHERE id = $3
`

type UpdateUserAllowIPParams struct {
	AllowedIps   []string  `json:"allowed_ips"`
	AllowipToken string    `json:"allowip_token"`
	ID           uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserAllowIP(ctx context.Context, arg UpdateUserAllowIPParams) error {
	_, err := q.db.ExecContext(ctx, updateUserAllowIP, pq.Array(arg.AllowedIps), arg.AllowipToken, arg.ID)
	return err
}

const updateUserAllowIPToken = `-- name: UpdateUserAllowIPToken :exec
UPDATE users SET 
    allowip_token = $1
WHERE id = $2
`

type UpdateUserAllowIPTokenParams struct {
	AllowipToken string    `json:"allowip_token"`
	ID           uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserAllowIPToken(ctx context.Context, arg UpdateUserAllowIPTokenParams) error {
	_, err := q.db.ExecContext(ctx, updateUserAllowIPToken, arg.AllowipToken, arg.ID)
	return err
}

const updateUserDeleteToken = `-- name: UpdateUserDeleteToken :exec
UPDATE users SET 
    delete_token = $1,
    last_delete_sent_at = $2
WHERE id = $3
`

type UpdateUserDeleteTokenParams struct {
	DeleteToken      string    `json:"delete_token"`
	LastDeleteSentAt time.Time `json:"last_delete_sent_at"`
	ID               uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserDeleteToken(ctx context.Context, arg UpdateUserDeleteTokenParams) error {
	_, err := q.db.ExecContext(ctx, updateUserDeleteToken, arg.DeleteToken, arg.LastDeleteSentAt, arg.ID)
	return err
}

const updateUserEmailChangeToken = `-- name: UpdateUserEmailChangeToken :exec
UPDATE users SET 
    emailchange_token = $1, 
    last_emailchange_sent_at = $2 
WHERE id = $3
`

type UpdateUserEmailChangeTokenParams struct {
	EmailchangeToken      string    `json:"emailchange_token"`
	LastEmailchangeSentAt time.Time `json:"last_emailchange_sent_at"`
	ID                    uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserEmailChangeToken(ctx context.Context, arg UpdateUserEmailChangeTokenParams) error {
	_, err := q.db.ExecContext(ctx, updateUserEmailChangeToken, arg.EmailchangeToken, arg.LastEmailchangeSentAt, arg.ID)
	return err
}

const updateUserEmailVerified = `-- name: UpdateUserEmailVerified :exec
UPDATE users SET 
    email_verified = $1,
    verify_token = $2
WHERE id = $3
`

type UpdateUserEmailVerifiedParams struct {
	EmailVerified bool      `json:"email_verified"`
	VerifyToken   string    `json:"verify_token"`
	ID            uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserEmailVerified(ctx context.Context, arg UpdateUserEmailVerifiedParams) error {
	_, err := q.db.ExecContext(ctx, updateUserEmailVerified, arg.EmailVerified, arg.VerifyToken, arg.ID)
	return err
}

const updateUserRecoveryToken = `-- name: UpdateUserRecoveryToken :exec
UPDATE users SET 
    recovery_token = $1,
    last_recovery_sent_at = $2
WHERE id = $3
`

type UpdateUserRecoveryTokenParams struct {
	RecoveryToken      string    `json:"recovery_token"`
	LastRecoverySentAt time.Time `json:"last_recovery_sent_at"`
	ID                 uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserRecoveryToken(ctx context.Context, arg UpdateUserRecoveryTokenParams) error {
	_, err := q.db.ExecContext(ctx, updateUserRecoveryToken, arg.RecoveryToken, arg.LastRecoverySentAt, arg.ID)
	return err
}

const updateUserResetPassword = `-- name: UpdateUserResetPassword :exec
UPDATE users SET 
    password = $1,
    recovery_token = $2
WHERE id = $3
`

type UpdateUserResetPasswordParams struct {
	Password      string    `json:"password"`
	RecoveryToken string    `json:"recovery_token"`
	ID            uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserResetPassword(ctx context.Context, arg UpdateUserResetPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserResetPassword, arg.Password, arg.RecoveryToken, arg.ID)
	return err
}

const updateUserVerifyToken = `-- name: UpdateUserVerifyToken :exec
UPDATE users SET 
    verify_token = $1,
    last_verify_sent_at = $2
WHERE id = $3
`

type UpdateUserVerifyTokenParams struct {
	VerifyToken      string    `json:"verify_token"`
	LastVerifySentAt time.Time `json:"last_verify_sent_at"`
	ID               uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserVerifyToken(ctx context.Context, arg UpdateUserVerifyTokenParams) error {
	_, err := q.db.ExecContext(ctx, updateUserVerifyToken, arg.VerifyToken, arg.LastVerifySentAt, arg.ID)
	return err
}