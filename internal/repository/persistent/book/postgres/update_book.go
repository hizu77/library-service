package book

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/pkg/transactor"
	"github.com/jackc/pgx/v5"
)

func (r *RepositoryImpl) UpdateBook(ctx context.Context, book entity.Book) (outBook entity.Book, txErr error) {
	var (
		tx  pgx.Tx
		err error
	)

	if tx, err = transactor.ExtractTx(ctx); err != nil {
		tx, err = r.Pool.Begin(ctx)
		if err != nil {
			return entity.Book{}, err
		}

		defer func() {
			if txErr != nil {
				tx.Rollback(ctx)
			} else {
				tx.Commit(ctx)
			}
		}()
	}

	sql, args, err := r.Builder.
		Update(TableName).
		Set(Name, book.Name).
		Where(squirrel.Eq{ID: book.ID}).
		ToSql()

	if err != nil {
		return entity.Book{}, err
	}

	cmdTag, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return entity.Book{}, err
	}

	if cmdTag.RowsAffected() == 0 {
		return entity.Book{}, entity.ErrBookNotFound
	}

	sql, args, err = r.Builder.
		Delete(AuthorBookTableName).
		Where(squirrel.Eq{BookID: book.ID}).
		ToSql()

	if err != nil {
		return entity.Book{}, err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return entity.Book{}, err
	}

	builder := r.Builder.Insert(AuthorBookTableName).Columns(AuthorID, BookID)

	for _, authorID := range book.AuthorsIDs {
		builder = builder.Values(authorID, book.ID)
	}

	sql, args, err = builder.ToSql()
	if err != nil {
		return entity.Book{}, err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return entity.Book{}, err
	}

	return book, nil
}
