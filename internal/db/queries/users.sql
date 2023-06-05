-- name: CreateUser :one
INSERT INTO users (id, email, username, password, email_verified, firstname, lastname, allowed_ips) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;


-- name: ReadUserByIdentity :one
SELECT * FROM users WHERE email = $1 OR username = $2 LIMIT 1;

-- name: ReadUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: ReadUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ReadUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: ReadUsers :many
SELECT * FROM users LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: UpdateUser :one
UPDATE users SET
    email = $1,
    username = $2,
    password = $3,
    firstname = $4,
    lastname = $5
WHERE id = $6 RETURNING *;


-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

