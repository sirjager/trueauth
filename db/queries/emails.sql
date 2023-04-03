-- name: CreateEmail :one
INSERT INTO emails (email,verified,token,last_token_sent_at) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetEmailByID :one
SELECT * FROM emails WHERE id = $1 LIMIT 1;


-- name: GetEmailByEmail :one
SELECT * FROM emails WHERE email = $1 LIMIT 1;


-- name: UpdateEmail :one
UPDATE emails SET
 verified = $1,
 token = $2,
 last_token_sent_at = $3
WHERE id = $4 RETURNING *;


