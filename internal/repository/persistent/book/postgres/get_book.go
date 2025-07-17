package book

import (
	"context"
	db "database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/hizu77/library-service/internal/entity"
)

func (r *RepositoryImpl) GetBook(ctx context.Context, id string) (entity.Book, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return entity.Book{}, err
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Select(Name).
		From(TableName).
		Where(squirrel.Eq{ID: id}).
		ToSql()

	if err != nil {
		return entity.Book{}, err
	}

	book := entity.Book{
		ID: id,
	}
	err = tx.QueryRow(ctx, sql, args...).Scan(&book.Name)

	if errors.Is(err, db.ErrNoRows) {
		return entity.Book{}, entity.ErrBookNotFound
	}

	if err != nil {
		return entity.Book{}, err
	}

	sql, args, err = r.Builder.
		Select(AuthorID).
		From(AuthorBookTableName).
		Where(squirrel.Eq{BookID: id}).
		ToSql()

	if err != nil {
		return entity.Book{}, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return entity.Book{}, err
	}
	defer rows.Close()

	var authorID string
	for rows.Next() {
		if err = rows.Scan(&authorID); err != nil {
			return entity.Book{}, err
		}

		book.AuthorsIDs = append(book.AuthorsIDs, authorID)
	}

	if err = tx.Commit(ctx); err != nil {
		return entity.Book{}, err
	}

	return book, nil
}
