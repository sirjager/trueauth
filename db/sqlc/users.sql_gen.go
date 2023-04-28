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

const create_User = `-- name: Create_User :one
INSERT INTO users (
    email, username, password,firstname, lastname,allowed_ips,
    last_verify_sent_at,last_recovery_sent_at,last_emailchange_sent_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at
`

type Create_UserParams struct {
	Email                 string    `json:"email"`
	Username              string    `json:"username"`
	Password              string    `json:"password"`
	Firstname             string    `json:"firstname"`
	Lastname              string    `json:"lastname"`
	AllowedIps            []string  `json:"allowed_ips"`
	LastVerifySentAt      time.Time `json:"last_verify_sent_at"`
	LastRecoverySentAt    time.Time `json:"last_recovery_sent_at"`
	LastEmailchangeSentAt time.Time `json:"last_emailchange_sent_at"`
}

func (q *Queries) Create_User(ctx context.Context, arg Create_UserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, create_User,
		arg.Email,
		arg.Username,
		arg.Password,
		arg.Firstname,
		arg.Lastname,
		pq.Array(arg.AllowedIps),
		arg.LastVerifySentAt,
		arg.LastRecoverySentAt,
		arg.LastEmailchangeSentAt,
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

const delete_User = `-- name: Delete_User :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) Delete_User(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, delete_User, id)
	return err
}

const read_User_ByEmail = `-- name: Read_User_ByEmail :one
SELECT id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at FROM users WHERE email = $1 LIMIT 1
`

func (q *Queries) Read_User_ByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, read_User_ByEmail, email)
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

const read_User_ByID = `-- name: Read_User_ByID :one
SELECT id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at FROM users WHERE id = $1 LIMIT 1
`

func (q *Queries) Read_User_ByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, read_User_ByID, id)
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

const read_User_ByUsername = `-- name: Read_User_ByUsername :one
SELECT id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at FROM users WHERE username = $1 LIMIT 1
`

func (q *Queries) Read_User_ByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, read_User_ByUsername, username)
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

const read_Users = `-- name: Read_Users :many
SELECT id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at FROM users LIMIT $2 OFFSET $1
`

type Read_UsersParams struct {
	Offset sql.NullInt32 `json:"offset"`
	Limit  sql.NullInt32 `json:"limit"`
}

func (q *Queries) Read_Users(ctx context.Context, arg Read_UsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, read_Users, arg.Offset, arg.Limit)
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

const update_User = `-- name: Update_User :one
UPDATE users SET
    username = $1,
    password = $2,
    firstname = $3,
    lastname = $4
    WHERE id = $5
RETURNING id, email, username, password, firstname, lastname, email_verified, verify_token, last_verify_sent_at, recovery_token, last_recovery_sent_at, emailchange_token, last_emailchange_sent_at, allowed_ips, allowip_token, delete_token, last_delete_sent_at, created_at, updated_at
`

type Update_UserParams struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) Update_User(ctx context.Context, arg Update_UserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, update_User,
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

const update_User_AllowIP = `-- name: Update_User_AllowIP :exec
UPDATE users SET 
    allowed_ips = $1,
    allowip_token = $2
WHERE id = $3
`

type Update_User_AllowIPParams struct {
	AllowedIps   []string  `json:"allowed_ips"`
	AllowipToken string    `json:"allowip_token"`
	ID           uuid.UUID `json:"id"`
}

func (q *Queries) Update_User_AllowIP(ctx context.Context, arg Update_User_AllowIPParams) error {
	_, err := q.db.ExecContext(ctx, update_User_AllowIP, pq.Array(arg.AllowedIps), arg.AllowipToken, arg.ID)
	return err
}

const update_User_AllowIPToken = `-- name: Update_User_AllowIPToken :exec
UPDATE users SET 
    allowip_token = $1
WHERE id = $2
`

type Update_User_AllowIPTokenParams struct {
	AllowipToken string    `json:"allowip_token"`
	ID           uuid.UUID `json:"id"`
}

func (q *Queries) Update_User_AllowIPToken(ctx context.Context, arg Update_User_AllowIPTokenParams) error {
	_, err := q.db.ExecContext(ctx, update_User_AllowIPToken, arg.AllowipToken, arg.ID)
	return err
}

const update_User_DeleteToken = `-- name: Update_User_DeleteToken :exec
UPDATE users SET 
    delete_token = $1,
    last_delete_sent_at = $2
WHERE id = $3
`

type Update_User_DeleteTokenParams struct {
	DeleteToken      string    `json:"delete_token"`
	LastDeleteSentAt time.Time `json:"last_delete_sent_at"`
	ID               uuid.UUID `json:"id"`
}

func (q *Queries) Update_User_DeleteToken(ctx context.Context, arg Update_User_DeleteTokenParams) error {
	_, err := q.db.ExecContext(ctx, update_User_DeleteToken, arg.DeleteToken, arg.LastDeleteSentAt, arg.ID)
	return err
}

const update_User_EmailChangeToken = `-- name: Update_User_EmailChangeToken :exec
UPDATE users SET 
    emailchange_token = $1, 
    last_emailchange_sent_at = $2 
WHERE id = $3
`

type Update_User_EmailChangeTokenParams struct {
	EmailchangeToken      string    `json:"emailchange_token"`
	LastEmailchangeSentAt time.Time `json:"last_emailchange_sent_at"`
	ID                    uuid.UUID `json:"id"`
}

func (q *Queries) Update_User_EmailChangeToken(ctx context.Context, arg Update_User_EmailChangeTokenParams) error {
	_, err := q.db.ExecContext(ctx, update_User_EmailChangeToken, arg.EmailchangeToken, arg.LastEmailchangeSentAt, arg.ID)
	return err
}

const update_User_EmailVerified = `-- name: Update_User_EmailVerified :exec
UPDATE users SET 
    email_verified = $1,
    verify_token = $2
WHERE id = $3
`

type Update_User_EmailVerifiedParams struct {
	EmailVerified bool      `json:"email_verified"`
	VerifyToken   string    `json:"verify_token"`
	ID            uuid.UUID `json:"id"`
}

func (q *Queries) Update_User_EmailVerified(ctx context.Context, arg Update_User_EmailVerifiedParams) error {
	_, err := q.db.ExecContext(ctx, update_User_EmailVerified, arg.EmailVerified, arg.VerifyToken, arg.ID)
	return err
}

const update_User_RecoveryToken = `-- name: Update_User_RecoveryToken :exec
UPDATE users SET 
    recovery_token = $1,
    last_recovery_sent_at = $2
WHERE id = $3
`

type Update_User_RecoveryTokenParams struct {
	RecoveryToken      string    `json:"recovery_token"`
	LastRecoverySentAt time.Time `json:"last_recovery_sent_at"`
	ID                 uuid.UUID `json:"id"`
}

func (q *Queries) Update_User_RecoveryToken(ctx context.Context, arg Update_User_RecoveryTokenParams) error {
	_, err := q.db.ExecContext(ctx, update_User_RecoveryToken, arg.RecoveryToken, arg.LastRecoverySentAt, arg.ID)
	return err
}

const update_User_ResetPassword = `-- name: Update_User_ResetPassword :exec
UPDATE users SET 
    password = $1,
    recovery_token = $2
WHERE id = $3
`

type Update_User_ResetPasswordParams struct {
	Password      string    `json:"password"`
	RecoveryToken string    `json:"recovery_token"`
	ID            uuid.UUID `json:"id"`
}

func (q *Queries) Update_User_ResetPassword(ctx context.Context, arg Update_User_ResetPasswordParams) error {
	_, err := q.db.ExecContext(ctx, update_User_ResetPassword, arg.Password, arg.RecoveryToken, arg.ID)
	return err
}

const update_User_VerifyToken = `-- name: Update_User_VerifyToken :exec
UPDATE users SET 
    verify_token = $1,
    last_verify_sent_at = $2
WHERE id = $3
`

type Update_User_VerifyTokenParams struct {
	VerifyToken      string    `json:"verify_token"`
	LastVerifySentAt time.Time `json:"last_verify_sent_at"`
	ID               uuid.UUID `json:"id"`
}

func (q *Queries) Update_User_VerifyToken(ctx context.Context, arg Update_User_VerifyTokenParams) error {
	_, err := q.db.ExecContext(ctx, update_User_VerifyToken, arg.VerifyToken, arg.LastVerifySentAt, arg.ID)
	return err
}
