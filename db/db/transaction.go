package db

import (
	"context"
	"fmt"
)

func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	queries := New(tx)
	if err = fn(queries); err != nil {
		if rollBackErr := tx.Rollback(ctx); rollBackErr != nil {
			return fmt.Errorf("transaction error: %v , rollback error: %v", err, rollBackErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
