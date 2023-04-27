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

type UpdateAccountEmailConfirmationTokenTxParams struct {
	UpdateAccountEmailConfirmationTokenParams
	BeforeUpdate func() error
}

func (store *SQLStore) UpdateAccountEmailConfirmationTokenTx(ctx context.Context, arg UpdateAccountEmailConfirmationTokenTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		if err = q.UpdateAccountEmailConfirmationToken(ctx, arg.UpdateAccountEmailConfirmationTokenParams); err != nil {
			return err
		}
		return
	})
}

type UpdateAccountEmailVerifiedTxParams struct {
	UpdateAccountEmailVerifiedParams
	AfterUpdate func() error
}

func (store *SQLStore) UpdateAccountEmailVerifiedTx(ctx context.Context, arg UpdateAccountEmailVerifiedTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if err = q.UpdateAccountEmailVerified(ctx, arg.UpdateAccountEmailVerifiedParams); err != nil {
			return err
		}
		if arg.AfterUpdate != nil {
			if err = arg.AfterUpdate(); err != nil {
				return err
			}
		}
		return
	})
}
