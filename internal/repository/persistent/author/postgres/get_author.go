package author

import (
	"context"
	db "database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/hizu77/library-service/internal/entity"
)

func (r *RepositoryImpl) GetAuthor(ctx context.Context, id string) (entity.Author, error) {
	sql, args, err := r.Builder.
		Select(Name).
		From(TableName).
		Where(squirrel.Eq{ID: id}).
		ToSql()

	if err != nil {
		return entity.Author{}, err
	}

	author := entity.Author{
		ID: id,
	}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&author.Name)

	if errors.Is(err, db.ErrNoRows) {
		return entity.Author{}, entity.ErrAuthorNotFound
	}

	if err != nil {
		return entity.Author{}, err
	}

	return author, nil
}
