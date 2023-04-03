// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: iprecords.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createIPRecord = `-- name: CreateIPRecord :exec
INSERT INTO iprecords (user_id, allowed_ips, blocked_ips, token) VALUES ($1, $2, $3, $4)
`

type CreateIPRecordParams struct {
	UserID     uuid.UUID `json:"user_id"`
	AllowedIps []string  `json:"allowed_ips"`
	BlockedIps []string  `json:"blocked_ips"`
	Token      string    `json:"token"`
}

func (q *Queries) CreateIPRecord(ctx context.Context, arg CreateIPRecordParams) error {
	_, err := q.db.ExecContext(ctx, createIPRecord,
		arg.UserID,
		pq.Array(arg.AllowedIps),
		pq.Array(arg.BlockedIps),
		arg.Token,
	)
	return err
}

const getIPRecordByID = `-- name: GetIPRecordByID :one
SELECT id, user_id, allowed_ips, blocked_ips, token, created_at, updated_at FROM iprecords WHERE id = $1 LIMIT 1
`

func (q *Queries) GetIPRecordByID(ctx context.Context, id uuid.UUID) (Iprecord, error) {
	row := q.db.QueryRowContext(ctx, getIPRecordByID, id)
	var i Iprecord
	err := row.Scan(
		&i.ID,
		&i.UserID,
		pq.Array(&i.AllowedIps),
		pq.Array(&i.BlockedIps),
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getIPRecordByUserID = `-- name: GetIPRecordByUserID :one
SELECT id, user_id, allowed_ips, blocked_ips, token, created_at, updated_at FROM iprecords WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetIPRecordByUserID(ctx context.Context, userID uuid.UUID) (Iprecord, error) {
	row := q.db.QueryRowContext(ctx, getIPRecordByUserID, userID)
	var i Iprecord
	err := row.Scan(
		&i.ID,
		&i.UserID,
		pq.Array(&i.AllowedIps),
		pq.Array(&i.BlockedIps),
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateIPRecord = `-- name: UpdateIPRecord :one
UPDATE iprecords SET
 allowed_ips = $1,
 blocked_ips = $2,
 token = $3
WHERE id = $4 RETURNING id, user_id, allowed_ips, blocked_ips, token, created_at, updated_at
`

type UpdateIPRecordParams struct {
	AllowedIps []string  `json:"allowed_ips"`
	BlockedIps []string  `json:"blocked_ips"`
	Token      string    `json:"token"`
	ID         uuid.UUID `json:"id"`
}

func (q *Queries) UpdateIPRecord(ctx context.Context, arg UpdateIPRecordParams) (Iprecord, error) {
	row := q.db.QueryRowContext(ctx, updateIPRecord,
		pq.Array(arg.AllowedIps),
		pq.Array(arg.BlockedIps),
		arg.Token,
		arg.ID,
	)
	var i Iprecord
	err := row.Scan(
		&i.ID,
		&i.UserID,
		pq.Array(&i.AllowedIps),
		pq.Array(&i.BlockedIps),
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateIPRecordTokenByUserID = `-- name: UpdateIPRecordTokenByUserID :exec
UPDATE iprecords SET token = $1 WHERE user_id = $2
`

type UpdateIPRecordTokenByUserIDParams struct {
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) UpdateIPRecordTokenByUserID(ctx context.Context, arg UpdateIPRecordTokenByUserIDParams) error {
	_, err := q.db.ExecContext(ctx, updateIPRecordTokenByUserID, arg.Token, arg.UserID)
	return err
}
