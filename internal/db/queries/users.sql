-- name: CreateUser :one
INSERT INTO users (
    email, username, password,firstname, lastname,allowed_ips,
    last_verify_sent_at,last_recovery_sent_at,last_emailchange_sent_at,last_delete_sent_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;


-- name: ReadUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: ReadUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ReadUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: ReadUsers :many
SELECT * FROM users LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');


-- name: UpdateUserDeleteToken :exec
UPDATE users SET 
    delete_token = $1,
    last_delete_sent_at = $2
WHERE id = $3;


-- name: UpdateUserVerifyToken :exec
UPDATE users SET 
    verify_token = $1,
    last_verify_sent_at = $2
WHERE id = $3;

-- name: UpdateUserEmailVerified :exec
UPDATE users SET 
    email_verified = $1,
    verify_token = $2
WHERE id = $3;


-- name: UpdateUserRecoveryToken :exec
UPDATE users SET 
    recovery_token = $1,
    last_recovery_sent_at = $2
WHERE id = $3;

-- name: UpdateUserResetPassword :exec
UPDATE users SET 
    password = $1,
    recovery_token = $2
WHERE id = $3;


-- name: UpdateUserEmailChangeToken :exec
UPDATE users SET 
    emailchange_token = $1, 
    last_emailchange_sent_at = $2 
WHERE id = $3;


-- name: UpdateUserAllowIPToken :exec
UPDATE users SET 
    allowip_token = $1
WHERE id = $2;


-- name: UpdateUserAllowIP :exec
UPDATE users SET 
    allowed_ips = $1,
    allowip_token = $2
WHERE id = $3;


-- name: UpdateUser :one
UPDATE users SET
    username = $1,
    password = $2,
    firstname = $3,
    lastname = $4
    WHERE id = $5
RETURNING *;


-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

