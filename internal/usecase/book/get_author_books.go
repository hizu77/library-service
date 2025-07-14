package book

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) GetAuthorBooks(ctx context.Context, id string) ([]entity.Book, error) {
	if _, err := u.authorRepository.GetAuthor(ctx, id); err != nil {
		u.logger.Error("authorRepository.GetAuthor", zap.String("id", id))
		return nil, entity.ErrAuthorNotFound
	}

	books, err := u.bookRepository.GetBooksByAuthorID(ctx, id)
	if err != nil {
		u.logger.Error("bookRepository.GetBooksByAuthor", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	u.logger.Info("GetAuthorBooks", zap.String("ID", id))

	return books, nil
}
