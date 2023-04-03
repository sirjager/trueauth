package sqlc

import (
	"context"

	"github.com/google/uuid"
)

type UpdateIPRecordTxParams struct {
	UpdateIPRecordParams
	BeforeUpdate func() error
	AfterUpdate  func(updated Iprecord) error
}

func (store *SQLStore) UpdateIPRecordTx(ctx context.Context, arg UpdateIPRecordTxParams) (Iprecord, error) {
	var ipRecord Iprecord
	err := store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}

		ipRecord, err = q.UpdateIPRecord(ctx, arg.UpdateIPRecordParams)
		if err != nil {
			return err
		}

		if arg.AfterUpdate != nil {
			err := arg.AfterUpdate(ipRecord)
			if err != nil {
				return err
			}
		}
		return
	})
	return ipRecord, err
}

type UpdateIPRecordTokenTxParams struct {
	UserID       uuid.UUID
	Token        string
	BeforeUpdate func() error
}

func (store *SQLStore) UpdateIPRecordTokenTx(ctx context.Context, arg UpdateIPRecordTokenTxParams) (Iprecord, error) {
	var ipRecord Iprecord
	err := store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		return q.UpdateIPRecordTokenByUserID(ctx, UpdateIPRecordTokenByUserIDParams{UserID: arg.UserID, Token: arg.Token})
	})
	return ipRecord, err
}
