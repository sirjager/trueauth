package sqlc

import "context"

type QuerierExtended interface {
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (User, error)
}
