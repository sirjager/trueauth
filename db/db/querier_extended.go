package db

import "context"

type QuerierExtended interface {
	UpdateUserPasswordTx(ctx context.Context, arg UpdateUserPasswordTx) error
}
