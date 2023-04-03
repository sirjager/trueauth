package sqlc

import "context"

type QuerierExtended interface {
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (User, error)

	UpdateIPRecordTx(ctx context.Context, arg UpdateIPRecordTxParams) (Iprecord, error)
	UpdateIPRecordTokenTx(ctx context.Context, arg UpdateIPRecordTokenTxParams) (Iprecord, error)

	CreateEmailRecordTx(ctx context.Context, arg CreateEmailRecordTxParams) (Emailrecord, error)
	UpdateEmailRecordTx(ctx context.Context, arg UpdateEmailRecordTxParams) (Emailrecord, error)
}
