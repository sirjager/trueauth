-- name: CreateAccount :one
INSERT INTO accounts (email,username,password,firstname,lastname) VALUES ($1, $2, $3, $4, $5) RETURNING *;


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


-- name: UpdateAccount :one
UPDATE accounts SET
 firstname = coalesce(sqlc.narg('firstname'), firstname),
 lastname = coalesce(sqlc.narg('lastname'), lastname),
 username = coalesce(sqlc.narg('username'), username),
 password = coalesce(sqlc.narg('password'), password)
WHERE id = $1 RETURNING *;


