package sqlc

import (
	"context"

	"github.com/google/uuid"
)

type UpdateIPTxParams struct {
	UpdateIPParams
	BeforeUpdate func(params UpdateIPParams) (UpdateIPParams, error)
	AfterUpdate  func(ip Ip) error
}

func (store *SQLStore) UpdateIPTx(ctx context.Context, arg UpdateIPTxParams) (Ip, error) {
	var ip Ip
	err := store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if arg.UpdateIPParams, err = arg.BeforeUpdate(arg.UpdateIPParams); err != nil {
				return err
			}
		}

		ip, err = q.UpdateIP(ctx, arg.UpdateIPParams)
		if err != nil {
			return err
		}

		if arg.AfterUpdate != nil {
			err := arg.AfterUpdate(ip)
			if err != nil {
				return err
			}
		}
		return
	})
	return ip, err
}

type UpdateIPTokenTxParams struct {
	AccountID    uuid.UUID
	Token        string
	BeforeUpdate func() error
}

func (store *SQLStore) UpdateIPTokenTx(ctx context.Context, arg UpdateIPTokenTxParams) (Ip, error) {
	var ip Ip
	err := store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}

		ip, err = q.UpdateIPTokenByAccountID(ctx, UpdateIPTokenByAccountIDParams{
			AccountID: arg.AccountID,
			Token:     arg.Token,
		})

		return err
	})

	return ip, err
}
