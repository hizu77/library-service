package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/library-service/internal/infra/model/outbox"
	"github.com/hizu77/library-service/internal/infra/repository"
	"github.com/hizu77/library-service/pkg/postgres"
	"github.com/hizu77/library-service/pkg/transactor"
	"github.com/jackc/pgx/v5"
)

var _ repository.OutboxRepository = (*Impl)(nil)

type Impl struct {
	*postgres.Postgres
}

func (i *Impl) SendMessage(ctx context.Context, idKey string, kind outbox.Kind, msg []byte) (txErr error) {
	var (
		tx  pgx.Tx
		err error
	)

	if tx, err = transactor.ExtractTx(ctx); err != nil {
		tx, err = i.Pool.Begin(ctx)
		if err != nil {
			return err
		}

		defer func() {
			if txErr != nil {
				tx.Rollback(ctx)
			} else {
				tx.Commit(ctx)
			}
		}()
	}

	sql, args, err := i.Builder.
		Insert(TableName).
		Columns(IdempotencyKey, Kind, Data).
		Values(idKey, kind, msg).
		Suffix("ON CONFLICT (" + IdempotencyKey + ") DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (i *Impl) GetMessages(ctx context.Context, batchSize int, inProgressTTL time.Duration) (outData []outbox.Data, txErr error) {
	var (
		tx  pgx.Tx
		err error
	)

	if tx, err = transactor.ExtractTx(ctx); err != nil {
		tx, err = i.Pool.Begin(ctx)
		if err != nil {
			return nil, err
		}

		defer func() {
			if txErr != nil {
				tx.Rollback(ctx)
			} else {
				tx.Commit(ctx)
			}
		}()
	}

	//IDK HOW TO DO THIS QUERY USING SQUIRREL
	//TODO FIX USING SQUIRREL
	const sql = `
UPDATE outbox
SET status = 'IN_PROGRESS'
WHERE idempotency_key IN (
    SELECT idempotency_key
    FROM outbox
    WHERE
        (status = 'CREATED'
            OR (status = 'IN_PROGRESS' AND updated_at < now() - $1::interval))
    ORDER BY created_at
    LIMIT $2
    FOR UPDATE SKIP LOCKED
	)
	RETURNING idempotency_key, data, kind;`

	interval := fmt.Sprintf("%d ms", inProgressTTL.Milliseconds())

	rows, err := tx.Query(ctx, sql, interval, batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]outbox.Data, 0)
	for rows.Next() {
		var idempotencyKey string
		var kind outbox.Kind
		var data []byte

		if err := rows.Scan(&idempotencyKey, &kind, &data); err != nil {
			return nil, err
		}

		out = append(out, outbox.Data{
			IdempotencyKey: idempotencyKey,
			Kind:           kind,
			RawData:        data,
		})
	}

	return out, nil
}

func (i *Impl) MarkMessageAsProcessed(ctx context.Context, idKeys []string) (txErr error) {
	var (
		tx  pgx.Tx
		err error
	)

	if tx, err = transactor.ExtractTx(ctx); err != nil {
		tx, err = i.Pool.Begin(ctx)
		if err != nil {
			return err
		}

		defer func() {
			if txErr != nil {
				tx.Rollback(ctx)
			} else {
				tx.Commit(ctx)
			}
		}()
	}

	sql, args, err := i.Builder.
		Update(TableName).
		Set(Status, "SUCCESS").
		Where(squirrel.Eq{IdempotencyKey: idKeys}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func New(pg *postgres.Postgres) *Impl {
	return &Impl{pg}
}
