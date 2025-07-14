package book

import (
	"context"

	"github.com/google/uuid"
	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) AddBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	for i := range book.AuthorsIDs {
		_, err := u.authorRepository.GetAuthor(ctx, book.AuthorsIDs[i])

		if err != nil {
			u.logger.Error("authorRepository.GetAuthor", zap.Error(err))
			return entity.Book{}, err
		}
	}

	book.ID = uuid.New().String()

	v, err := u.bookRepository.AddBook(ctx, book)

	if err != nil {
		u.logger.Error("bookRepository.AddBook", zap.Error(err))
		return entity.Book{}, err
	}

	u.logger.Info("AddBook", zap.String("ID", v.ID))

	return v, nil
}
