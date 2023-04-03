package sqlc

import (
	"context"
)

type CreateAccountTxParams struct {
	CreateAccountParams
	AfterCreate func(account Account) error
}

func (store *SQLStore) CreateAccountTx(ctx context.Context, arg CreateAccountTxParams) (Account, error) {
	var user Account
	err := store.execTx(ctx, func(q *Queries) (err error) {
		user, err = q.CreateAccount(ctx, arg.CreateAccountParams)
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
