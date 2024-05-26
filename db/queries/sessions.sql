-- name: CreateSession :one
INSERT INTO "_sessions" (
  id, 
  refresh_token,
  access_token_id,
  access_token,
  user_id,
  client_ip,
  user_agent,
  blocked,
  access_token_expires_at,
  refresh_token_expires_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: ReadSession :one
SELECT * FROM "_sessions" WHERE id = $1 LIMIT 1;

-- name: ReadSessionByAccessTokenID :one
SELECT * FROM "_sessions" WHERE access_token_id = $1 LIMIT 1;

-- name: ReadSessionsByUserID :many
SELECT * FROM "_sessions" WHERE user_id = $1 LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: DeleteSessionByUserID :exec
DELETE FROM "_sessions" WHERE user_id = $1;

-- name: DeleteSession :exec
DELETE FROM "_sessions" WHERE id = $1;

-- name: DeleteSessionByAccessTokenID :exec
DELETE FROM "_sessions" WHERE access_token_id = $1;
