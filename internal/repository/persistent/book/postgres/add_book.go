package book

import (
	"context"
	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/repository/persistent/utils"
)

func (r *RepositoryImpl) AddBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return entity.Book{}, err
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Insert(TableName).
		Columns(ID, Name).
		Values(book.ID, book.Name).
		ToSql()

	if err != nil {
		return entity.Book{}, err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if utils.IsUniqueConstraintError(err) {
		return entity.Book{}, entity.ErrBookAlreadyExists
	}

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

	if err = tx.Commit(ctx); err != nil {
		return entity.Book{}, err
	}

	return book, nil
}
