package sqlc

import "context"

type QuerierExtended interface {
	GetUser(ctx context.Context, arg GetUserParams) (User, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (User, error)
}
