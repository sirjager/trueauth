-- name: Create_Session :one
INSERT INTO sessions (
    id, refresh_token, access_token_id, access_token, user_id, client_ip, user_agent,
    blocked, access_token_expires_at, refresh_token_expires_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: Read_Session_ByID :one
SELECT * FROM sessions WHERE id = $1 LIMIT 1;

-- name: Read_Session_ByAccessTokenID :one
SELECT * FROM sessions WHERE access_token_id = $1 LIMIT 1;

-- name: Read_Sessions_ByUserID :many
SELECT * FROM sessions WHERE user_id = $1 LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: Delete_Session :exec
DELETE FROM sessions WHERE id = $1;

-- name: Delete_Session_ByUserID :exec
DELETE FROM sessions WHERE user_id = $1;
