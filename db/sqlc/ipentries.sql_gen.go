// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: ipentries.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createIPEntry = `-- name: CreateIPEntry :exec
INSERT INTO ipentries (id, allowed_ips, blocked_ips, code, code_expires_at) VALUES ($1, $2, $3, $4, $5)
`

type CreateIPEntryParams struct {
	ID            uuid.UUID `json:"id"`
	AllowedIps    []string  `json:"allowed_ips"`
	BlockedIps    []string  `json:"blocked_ips"`
	Code          string    `json:"code"`
	CodeExpiresAt time.Time `json:"code_expires_at"`
}

func (q *Queries) CreateIPEntry(ctx context.Context, arg CreateIPEntryParams) error {
	_, err := q.db.ExecContext(ctx, createIPEntry,
		arg.ID,
		pq.Array(arg.AllowedIps),
		pq.Array(arg.BlockedIps),
		arg.Code,
		arg.CodeExpiresAt,
	)
	return err
}

const getIPEntry = `-- name: GetIPEntry :one
SELECT id, allowed_ips, blocked_ips, code, code_expires_at, created_at, updated_at FROM ipentries WHERE id = $1 LIMIT 1
`

func (q *Queries) GetIPEntry(ctx context.Context, id uuid.UUID) (Ipentry, error) {
	row := q.db.QueryRowContext(ctx, getIPEntry, id)
	var i Ipentry
	err := row.Scan(
		&i.ID,
		pq.Array(&i.AllowedIps),
		pq.Array(&i.BlockedIps),
		&i.Code,
		&i.CodeExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateIPEntry = `-- name: UpdateIPEntry :one
UPDATE ipentries SET
 allowed_ips = $1,
 blocked_ips = $2,
 code = $3,
 code_expires_at = $4
WHERE id = $5 RETURNING id, allowed_ips, blocked_ips, code, code_expires_at, created_at, updated_at
`

type UpdateIPEntryParams struct {
	AllowedIps    []string  `json:"allowed_ips"`
	BlockedIps    []string  `json:"blocked_ips"`
	Code          string    `json:"code"`
	CodeExpiresAt time.Time `json:"code_expires_at"`
	ID            uuid.UUID `json:"id"`
}

func (q *Queries) UpdateIPEntry(ctx context.Context, arg UpdateIPEntryParams) (Ipentry, error) {
	row := q.db.QueryRowContext(ctx, updateIPEntry,
		pq.Array(arg.AllowedIps),
		pq.Array(arg.BlockedIps),
		arg.Code,
		arg.CodeExpiresAt,
		arg.ID,
	)
	var i Ipentry
	err := row.Scan(
		&i.ID,
		pq.Array(&i.AllowedIps),
		pq.Array(&i.BlockedIps),
		&i.Code,
		&i.CodeExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
