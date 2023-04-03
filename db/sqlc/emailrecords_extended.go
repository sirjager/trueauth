package sqlc

import (
	"context"
)

type CreateEmailRecordTxParams struct {
	CreateEmailRecordParams
	BeforeCreate func() error
}

func (store *SQLStore) CreateEmailRecordTx(ctx context.Context, arg CreateEmailRecordTxParams) (Emailrecord, error) {
	var emailRecord Emailrecord
	err := store.execTx(ctx, func(q *Queries) (err error) {
		// send email before creating record
		if arg.BeforeCreate != nil {
			if err = arg.BeforeCreate(); err != nil {
				return err
			}
		}
		emailRecord, err = q.CreateEmailRecord(ctx, arg.CreateEmailRecordParams)
		if err != nil {
			return err
		}
		return
	})
	return emailRecord, err
}

type UpdateEmailRecordTxParams struct {
	UpdateEmailRecordParams
	BeforeUpdate func() error
	AfterUpdate  func(Emailrecord) error
}

func (store *SQLStore) UpdateEmailRecordTx(ctx context.Context, arg UpdateEmailRecordTxParams) (Emailrecord, error) {
	var emailRecord Emailrecord
	err := store.execTx(ctx, func(q *Queries) (err error) {

		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}

		emailRecord, err = q.UpdateEmailRecord(ctx, arg.UpdateEmailRecordParams)
		if err != nil {
			return err
		}

		if arg.AfterUpdate != nil {
			if err = arg.AfterUpdate(emailRecord); err != nil {
				return err
			}
		}

		return
	})
	return emailRecord, err
}
