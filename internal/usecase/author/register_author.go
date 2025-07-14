package author

import (
	"context"

	"github.com/google/uuid"
	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) RegisterAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	author.ID = uuid.New().String()

	v, err := u.authorRepository.AddAuthor(ctx, author)

	if err != nil {
		u.logger.Error("authorRepository.AddAuthor", zap.Error(err))
		return entity.Author{}, err
	}

	u.logger.Info("RegisterAuthor", zap.String("ID", v.ID))

	return v, nil
}
