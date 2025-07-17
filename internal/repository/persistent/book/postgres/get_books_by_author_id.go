package book

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/hizu77/library-service/internal/entity"
)

func (r *RepositoryImpl) GetBooksByAuthorID(ctx context.Context, authorID string) ([]entity.Book, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Select(
			"b."+ID,
			"b."+Name,
			"ARRAY_AGG(ab."+AuthorID+") AS author_ids",
		).
		From(TableName+" b").
		Join(AuthorBookTableName+" ab ON b."+ID+" = ab."+BookID).
		Where(squirrel.Eq{"ab." + AuthorID: authorID}).
		GroupBy("b."+ID, "b."+Name).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []entity.Book
	var book entity.Book
	for rows.Next() {
		if err = rows.Scan(&book.ID, &book.Name, &book.AuthorsIDs); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return books, nil
}
