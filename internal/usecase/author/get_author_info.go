package author

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) GetAuthorInfo(ctx context.Context, id string) (entity.Author, error) {
	author, err := u.authorRepository.GetAuthor(ctx, id)

	if err != nil {
		u.logger.Error("authorRepository.GetAuthor", zap.Error(err))
		return entity.Author{}, err
	}

	u.logger.Info("GetAuthorInfo", zap.String("ID", author.ID))

	return author, nil
}
