-- name: CreateSession :one
INSERT INTO sessions (
    id, refresh_token, access_token_id, access_token, user_id,
    blocked, access_token_expires_at, refresh_token_expires_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions WHERE id = $1 LIMIT 1;

-- name: ListSessionsByUser :many
SELECT * FROM sessions WHERE user_id = $1 LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

