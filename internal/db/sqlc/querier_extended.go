package sqlc

import "context"

type QuerierExtended interface {
	UpdateUserVerifyTokenTx(ctx context.Context, arg UpdateUserVerifyTokenTxParams) error
	UpdateUserEmailVerifiedTx(ctx context.Context, arg UpdateUserEmailVerifiedTxParams) error

	UpdateUserRecoveryTokenTx(ctx context.Context, arg UpdateUserRecoveryTokenTxParams) error
	UpdateUserResetPasswordTx(ctx context.Context, arg UpdateUserResetPasswordTxParams) error

	UpdateUserDeleteTokenTx(ctx context.Context, arg UpdateUserDeleteTokenTxParams) error
}
