package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

const userTableFields = " id, email, username, password, firstname, lastname, verified, blocked, created_at, updated_at "

const getUser = "SELECT" + userTableFields + "FROM users WHERE $1 = $2 LIMIT 1 "

func _scanUser(row *sql.Row, i User) error {
	return row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Firstname,
		&i.Lastname,
		&i.Verified,
		&i.Blocked,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
}

type GetUserParams struct {
	Identity string `json:"identity"`
	Is       string `json:"is"`
}

func (q *Queries) GetUser(ctx context.Context, arg GetUserParams) (User, error) {
	var i User
	valid := false
	if arg.Identity == "id" || arg.Identity == "email" || arg.Identity == "username" {
		valid = true
	}
	if !valid {
		return i, fmt.Errorf("invalid identity: %s. identity can only be id, email or username", arg.Identity)
	}
	row := q.db.QueryRowContext(ctx, getUser, arg.Identity, arg.Is)
	err := _scanUser(row, i)
	return i, err
}
