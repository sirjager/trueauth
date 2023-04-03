-- name: CreateIP :exec
INSERT INTO ips (account_id, allowed_ips, blocked_ips, token) VALUES ($1, $2, $3, $4);

-- name: GetIPByID :one
SELECT * FROM ips WHERE id = $1 LIMIT 1;

-- name: GetIPByAccountID :one
SELECT * FROM ips WHERE account_id = $1 LIMIT 1;


-- name: UpdateIP :one
UPDATE ips SET
 allowed_ips = $1,
 blocked_ips = $2,
 token = $3
WHERE id = $4 RETURNING *;


-- name: UpdateIPTokenByAccountID :one
UPDATE ips SET token = $1 WHERE account_id = $2 RETURNING *;
