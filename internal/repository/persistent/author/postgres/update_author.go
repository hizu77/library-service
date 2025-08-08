package author

import (
	"context"
	db "database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/pkg/transactor"
	"github.com/jackc/pgx/v5"
)

func (r *RepositoryImpl) UpdateAuthor(ctx context.Context, author entity.Author) (outAuthor entity.Author, txErr error) {
	var (
		tx  pgx.Tx
		err error
	)

	if tx, err = transactor.ExtractTx(ctx); err != nil {
		tx, err = r.Pool.Begin(ctx)
		if err != nil {
			return entity.Author{}, err
		}

		defer func() {
			if txErr != nil {
				_ = tx.Rollback(ctx)
			} else {
				_ = tx.Commit(ctx)
			}
		}()
	}

	sql, args, err := r.Builder.
		Update(TableName).
		Set(Name, author.Name).
		Where(squirrel.Eq{ID: author.ID}).
		ToSql()
	if err != nil {
		return entity.Author{}, err
	}

	_, err = tx.Exec(ctx, sql, args...)

	if errors.Is(err, db.ErrNoRows) {
		return entity.Author{}, entity.ErrAuthorNotFound
	}

	if err != nil {
		return entity.Author{}, err
	}

	return author, nil
}
