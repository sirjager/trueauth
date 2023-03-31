-- name: CreateUser :one
INSERT INTO users (email,username,password,firstname,lastname) VALUES ($1, $2, $3, $4, $5) RETURNING *;


-- name: ListUsers :many
SELECT * FROM users LIMIT sqlc.narg('limit') OFFSET sqlc.narg('offset');

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users SET
 firstname = coalesce(sqlc.narg('firstname'), firstname),
 lastname = coalesce(sqlc.narg('lastname'), lastname),
 username = coalesce(sqlc.narg('username'), username),
 password = coalesce(sqlc.narg('password'), password)
WHERE id = $1 RETURNING *;


