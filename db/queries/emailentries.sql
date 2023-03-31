-- name: CreateEmailEntry :one
INSERT INTO emailentries (
    email, user_id, verified, code, code_expires_at
) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetEmailEntry :one
SELECT * FROM emailentries WHERE id = $1 LIMIT 1;

-- name: UpdateEmailEntry :one
UPDATE emailentries SET
 verified = $1,
 code = $2,
 code_expires_at = $3
WHERE id = $4 RETURNING *;


