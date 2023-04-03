-- name: CreateEmailRecord :one
INSERT INTO emailrecords (email,verified,token,last_token_sent_at) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetEmailRecordByID :one
SELECT * FROM emailrecords WHERE id = $1 LIMIT 1;


-- name: GetEmailRecordByEmail :one
SELECT * FROM emailrecords WHERE email = $1 LIMIT 1;


-- name: UpdateEmailRecord :one
UPDATE emailrecords SET
 verified = $1,
 token = $2,
 last_token_sent_at = $3
WHERE id = $4 RETURNING *;


