package db

import (
	"context"
)

type UpdateUserPasswordTx struct {
	BeforeUpdate func() error
	AfterUpdate  func() error
	UpdateUserPasswordParams
}

func (s *SQLStore) UpdateUserPasswordTx(ctx context.Context, arg UpdateUserPasswordTx) error {
	return s.execTx(ctx, func(q *Queries) error {
		if arg.BeforeUpdate != nil {
			if err := arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		err := q.UpdateUserPassword(ctx, arg.UpdateUserPasswordParams)
		if err != nil {
			return err
		}

		if arg.AfterUpdate != nil {
			if err := arg.AfterUpdate(); err != nil {
				return err
			}
		}
		return nil
	})
}
