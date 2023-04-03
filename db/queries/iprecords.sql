-- name: CreateIPRecord :exec
INSERT INTO iprecords (user_id, allowed_ips, blocked_ips, token) VALUES ($1, $2, $3, $4);

-- name: GetIPRecordByID :one
SELECT * FROM iprecords WHERE id = $1 LIMIT 1;

-- name: GetIPRecordByUserID :one
SELECT * FROM iprecords WHERE user_id = $1 LIMIT 1;


-- name: UpdateIPRecord :one
UPDATE iprecords SET
 allowed_ips = $1,
 blocked_ips = $2,
 token = $3
WHERE id = $4 RETURNING *;


-- name: UpdateIPRecordTokenByUserID :exec
UPDATE iprecords SET token = $1 WHERE user_id = $2;
