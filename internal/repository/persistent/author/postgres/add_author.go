package author

import (
	"context"
	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/repository/persistent/utils"
)

func (r *RepositoryImpl) AddAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return entity.Author{}, err
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Insert(TableName).
		Columns(ID, Name).
		Values(author.ID, author.Name).
		ToSql()

	if err != nil {
		return entity.Author{}, err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if utils.IsUniqueConstraintError(err) {
		return entity.Author{}, entity.ErrAuthorAlreadyExists
	}

	if err != nil {
		return entity.Author{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return entity.Author{}, err
	}

	return author, nil
}
