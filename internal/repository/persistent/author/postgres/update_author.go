package author

import (
	"context"
	db "database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/hizu77/library-service/internal/entity"
)

func (r *RepositoryImpl) UpdateAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return entity.Author{}, err
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Update(TableName).
		Set(Name, author.Name).
		Where(squirrel.Eq{Name: author.Name}).
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

	if err = tx.Commit(ctx); err != nil {
		return entity.Author{}, err
	}

	return author, nil
}
