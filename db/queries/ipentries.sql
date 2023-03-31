-- name: CreateIPEntry :exec
INSERT INTO ipentries (id, allowed_ips, blocked_ips, code, code_expires_at) VALUES ($1, $2, $3, $4, $5);

-- name: GetIPEntry :one
SELECT * FROM ipentries WHERE id = $1 LIMIT 1;

-- name: UpdateIPEntry :one
UPDATE ipentries SET
 allowed_ips = $1,
 blocked_ips = $2,
 code = $3,
 code_expires_at = $4
WHERE id = $5 RETURNING *;


