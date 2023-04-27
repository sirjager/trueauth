-- name: CreateAccount :one
INSERT INTO accounts (email,username,password,firstname,lastname,allowed_ips,last_confirmation_sent_at,last_recovery_sent_at,last_email_change_sent_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;


-- name: GetAccountByID :one
SELECT * FROM accounts WHERE id = $1 LIMIT 1;

-- name: GetAccountByEmail :one
SELECT * FROM accounts WHERE email = $1 LIMIT 1;

-- name: GetAccountByUsername :one
SELECT * FROM accounts WHERE username = $1 LIMIT 1;


-- name: ListAccounts :many
SELECT * FROM accounts LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;


-- name: UpdateAccountRecoveryToken :exec
UPDATE accounts SET 
recovery_token = $1, 
last_recovery_sent_at = $2 
WHERE id = $3;

-- name: UpdateAccountPassword :exec
UPDATE accounts SET 
password = $1 
WHERE id = $2;


-- name: UpdateAccountEmailChangeToken :exec
UPDATE accounts SET 
email_change_token = $1, 
last_email_change_sent_at = $2 
WHERE id = $3;

-- name: UpdateAccountEmail :exec
UPDATE accounts SET 
email = $1 
WHERE id = $2;


-- name: UpdateAccountAllowIPToken :exec
UPDATE accounts SET 
allow_ip_token = $1
WHERE id = $2;

-- name: UpdateAccountAllowIP :exec
UPDATE accounts SET 
allowed_ips = $1,
allow_ip_token = $2
WHERE id = $3;


-- name: UpdateAccountEmailConfirmationToken :exec
UPDATE accounts SET 
 confirmation_token = $1,
 last_confirmation_sent_at = $2
WHERE id = $3;

-- name: UpdateAccountEmailVerified :exec
UPDATE accounts SET 
 email_verified = $1,
 confirmation_token = $2
WHERE id = $3;
