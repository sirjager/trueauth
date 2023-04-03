package sqlc

import "context"

type QuerierExtended interface {
	CreateAccountTx(ctx context.Context, arg CreateAccountTxParams) (Account, error)

	CreateEmailTx(ctx context.Context, arg CreateEmailTxParams) (Email, error)
	UpdateEmailTx(ctx context.Context, arg UpdateEmailTxParams) (Email, error)

	UpdateIPTx(ctx context.Context, arg UpdateIPTxParams) (Ip, error)
	UpdateIPTokenTx(ctx context.Context, arg UpdateIPTokenTxParams) (Ip, error)
}
