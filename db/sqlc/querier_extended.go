package sqlc

import "context"

type QuerierExtended interface {
	GetUser(ctx context.Context, arg GetUserParams) (User, error)
}

var _ QuerierExtended = (*Queries)(nil)
