package book

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) GetBookInfo(ctx context.Context, id string) (entity.Book, error) {
	var outBook entity.Book

	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		var txErr error
		outBook, txErr = u.bookRepository.GetBook(ctx, id)
		if txErr != nil {
			u.logger.Error(
				"bookRepository.GetBook",
				zap.Error(txErr),
				zap.String("book_id", id),
			)

			return txErr
		}

		u.logger.Info(
			"GetBookInfo",
			zap.String("book_id", id),
		)

		return nil
	})
	if err != nil {
		return entity.Book{}, err
	}

	return outBook, nil
}
