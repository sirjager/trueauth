-- name: CreateEmailRecord :one
INSERT INTO emailrecords (
    email, user_id, verified, code, code_expires_at
) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetEmailRecord :one
SELECT * FROM emailrecords WHERE id = $1 LIMIT 1;

-- name: UpdateEmailRecord :one
UPDATE emailrecords SET
 verified = $1,
 code = $2,
 code_expires_at = $3
WHERE id = $4 RETURNING *;


