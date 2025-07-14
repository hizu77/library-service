package author

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) ChangeAuthorInfo(ctx context.Context, author entity.Author) (entity.Author, error) {
	v, err := u.authorRepository.UpdateAuthor(ctx, author)

	if err != nil {
		u.logger.Error("authorRepository.UpdateAuthor", zap.Error(err))
		return entity.Author{}, err
	}

	u.logger.Info("ChangeAuthorInfo", zap.String("ID", v.ID))

	return v, nil
}
