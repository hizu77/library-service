package book

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) UpdateBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	for i := range book.AuthorsIDs {
		_, err := u.authorRepository.GetAuthor(ctx, book.AuthorsIDs[i])
		if err != nil {
			u.logger.Error("authorRepository.GetAuthor", zap.Error(err))
			return entity.Book{}, err
		}
	}

	v, err := u.bookRepository.UpdateBook(ctx, book)

	if err != nil {
		u.logger.Error("bookRepository.UpdateBook", zap.Error(err))
		return entity.Book{}, err
	}

	u.logger.Info("UpdateBook", zap.String("ID", v.ID))

	return v, nil
}
