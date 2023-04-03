package sqlc

import (
	"context"
)

type CreateEmailTxParams struct {
	CreateEmailParams
	BeforeCreate func(params CreateEmailParams) (CreateEmailParams, error)
	AfterCreate  func(email Email) error
}

func (store *SQLStore) CreateEmailTx(ctx context.Context, arg CreateEmailTxParams) (Email, error) {
	var email Email
	err := store.execTx(ctx, func(q *Queries) (err error) {
		// send email before creating record
		if arg.BeforeCreate != nil {
			if arg.CreateEmailParams, err = arg.BeforeCreate(arg.CreateEmailParams); err != nil {
				return err
			}
		}

		email, err = q.CreateEmail(ctx, arg.CreateEmailParams)
		if err != nil {
			return err
		}

		if arg.AfterCreate != nil {
			if err = arg.AfterCreate(email); err != nil {
				return err
			}
		}

		return
	})
	return email, err
}

type UpdateEmailTxParams struct {
	UpdateEmailParams
	BeforeUpdate func() error
	AfterUpdate  func(Email) error
}

func (store *SQLStore) UpdateEmailTx(ctx context.Context, arg UpdateEmailTxParams) (Email, error) {
	var email Email
	err := store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}

		email, err = q.UpdateEmail(ctx, arg.UpdateEmailParams)
		if err != nil {
			return err
		}

		if arg.AfterUpdate != nil {
			if err = arg.AfterUpdate(email); err != nil {
				return err
			}
		}

		return
	})
	return email, err
}
