package sqlc

import (
	"context"
)

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (User, error) {
	var user User
	err := store.execTx(ctx, func(q *Queries) (err error) {
		user, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		if arg.AfterCreate != nil {
			if err = arg.AfterCreate(user); err != nil {
				return err
			}
		}

		return nil
	})

	return user, err
}
