-- name: ReadUsers :many
select * from "_users" limit sqlc.narg('limit') offset sqlc.narg('offset');

-- name: ReadUser :one
select * from "_users" where id = @id limit 1;

-- name: ReadUserByEmail :one
select * from "_users" where email = @email limit 1;

-- name: ReadUserByUsername :one
select * from "_users" where username = @username limit 1;

-- name: DeleteUser :exec
delete from "_users" where id = $1;

-- name: CreateUser :one
INSERT INTO "_users" (id,email,username,hash_salt,hash_pass,firstname,lastname,verified) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateUserEmailVerificationToken :exec
UPDATE "_users" SET token_email_verify = $1, last_email_verify = $2 WHERE id = @id;

-- name: UpdateUserEmailVerified :one
UPDATE "_users" SET verified = $1, token_email_verify = $2 WHERE id = @id RETURNING *;

-- name: UpdateUserDeletionToken :exec
UPDATE "_users" SET token_user_deletion = $1, last_user_deletion = $2 WHERE id = @id;

-- name: UpdateUserPasswordResetToken :exec
UPDATE "_users" SET token_password_reset = $1, last_password_reset = $2 WHERE id = @id;

-- name: UpdateUserPassword :exec 
UPDATE "_users" SET hash_pass = $1, hash_salt = $2, last_password_reset = $3 WHERE id = @id;

