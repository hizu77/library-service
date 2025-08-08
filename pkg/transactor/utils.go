package transactor

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	txInjector struct{}
)

var (
	ErrTxNotFound = errors.New("transaction not found")
)

func ExtractTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(txInjector{}).(pgx.Tx)
	if !ok {
		return nil, ErrTxNotFound
	}

	return tx, nil
}

func InjectTx(ctx context.Context, pool *pgxpool.Pool) (context.Context, pgx.Tx, error) {
	if tx, err := ExtractTx(ctx); err == nil {
		return ctx, tx, nil
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	return context.WithValue(ctx, txInjector{}, tx), tx, nil
}
