package author

import (
	"context"

	"github.com/google/uuid"
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

func (u *UseCaseImpl) ChangeAuthorInfo(ctx context.Context, author entity.Author) (entity.Author, error) {
	v, err := u.authorRepository.UpdateAuthor(ctx, author)

	if err != nil {
		u.logger.Error("authorRepository.UpdateAuthor", zap.Error(err))
		return entity.Author{}, err
	}

	u.logger.Info("ChangeAuthorInfo", zap.String("ID", v.ID))

	return v, nil
}
