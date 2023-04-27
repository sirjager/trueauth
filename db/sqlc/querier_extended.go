package sqlc

import "context"

type QuerierExtended interface {
	CreateAccountTx(ctx context.Context, arg CreateAccountTxParams) (Account, error)

	UpdateAccountEmailConfirmationTokenTx(ctx context.Context, arg UpdateAccountEmailConfirmationTokenTxParams) error
	UpdateAccountEmailVerifiedTx(ctx context.Context, arg UpdateAccountEmailVerifiedTxParams) error
}
