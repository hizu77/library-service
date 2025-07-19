package book

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) AddBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	var outBook entity.Book

	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		var txErr error

		for i := range book.AuthorsIDs {
			_, txErr = u.authorRepository.GetAuthor(ctx, book.AuthorsIDs[i])
			if txErr != nil {
				u.logger.Error("authorRepository.GetAuthor", zap.Error(txErr))
				return txErr
			}
		}

		outBook, txErr = u.bookRepository.AddBook(ctx, book)
		if txErr != nil {
			u.logger.Error("bookRepository.AddBook", zap.Error(txErr))
			return txErr
		}

		u.logger.Info("AddBook", zap.String("ID", outBook.ID))

		return nil
	})
	if err != nil {
		return entity.Book{}, err
	}

	return outBook, nil
}
