-- name: CreateIPRecord :exec
INSERT INTO iprecords (id, allowed_ips, blocked_ips, code, code_expires_at) VALUES ($1, $2, $3, $4, $5);

-- name: GetIPRecord :one
SELECT * FROM iprecords WHERE id = $1 LIMIT 1;

-- name: UpdateIPRecord :one
UPDATE iprecords SET
 allowed_ips = $1,
 blocked_ips = $2,
 code = $3,
 code_expires_at = $4
WHERE id = $5 RETURNING *;


