package transactor

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Transactor = (*Impl)(nil)

type Impl struct {
	pool *pgxpool.Pool
}

func (i *Impl) WithTx(ctx context.Context, function func(ctx context.Context) error) (txErr error) {
	ctxWithTx, tx, err := injectTx(ctx, i.pool)
	if err != nil {
		return fmt.Errorf("injectTx: %w", err)
	}

	defer func() {
		if txErr != nil {
			tx.Rollback(ctxWithTx)
			return
		}

		tx.Commit(ctxWithTx)
	}()

	err = function(ctx)
	if err != nil {
		return fmt.Errorf("function: %w", err)
	}

	return nil
}
