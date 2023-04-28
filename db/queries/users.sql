-- name: Create_User :one
INSERT INTO users (
    email, username, password,firstname, lastname,allowed_ips,
    last_verify_sent_at,last_recovery_sent_at,last_emailchange_sent_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;


-- name: Read_User_ByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: Read_User_ByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: Read_User_ByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: Read_Users :many
SELECT * FROM users LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');


-- name: Update_User_DeleteToken :exec
UPDATE users SET 
    delete_token = $1,
    last_delete_sent_at = $2
WHERE id = $3;


-- name: Update_User_VerifyToken :exec
UPDATE users SET 
    verify_token = $1,
    last_verify_sent_at = $2
WHERE id = $3;

-- name: Update_User_EmailVerified :exec
UPDATE users SET 
    email_verified = $1,
    verify_token = $2
WHERE id = $3;


-- name: Update_User_RecoveryToken :exec
UPDATE users SET 
    recovery_token = $1,
    last_recovery_sent_at = $2
WHERE id = $3;

-- name: Update_User_ResetPassword :exec
UPDATE users SET 
    password = $1,
    recovery_token = $2
WHERE id = $3;


-- name: Update_User_EmailChangeToken :exec
UPDATE users SET 
    emailchange_token = $1, 
    last_emailchange_sent_at = $2 
WHERE id = $3;


-- name: Update_User_AllowIPToken :exec
UPDATE users SET 
    allowip_token = $1
WHERE id = $2;


-- name: Update_User_AllowIP :exec
UPDATE users SET 
    allowed_ips = $1,
    allowip_token = $2
WHERE id = $3;


-- name: Update_User :one
UPDATE users SET
    username = $1,
    password = $2,
    firstname = $3,
    lastname = $4
    WHERE id = $5
RETURNING *;


-- name: Delete_User :exec
DELETE FROM users WHERE id = $1;

