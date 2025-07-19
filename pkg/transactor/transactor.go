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

func New(pool *pgxpool.Pool) *Impl {
	return &Impl{
		pool: pool,
	}
}

func (i *Impl) WithTx(ctx context.Context, function func(ctx context.Context) error) (txErr error) {
	ctxWithTx, tx, err := InjectTx(ctx, i.pool)
	if err != nil {
		return fmt.Errorf("InjectTx: %w", err)
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
