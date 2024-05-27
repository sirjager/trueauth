package db

import (
	_"github.com/golang/mock/mockgen/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	QuerierExtended
}

type SQLStore struct {
	*Queries
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		pool:    pool,
		Queries: New(pool),
	}
}
