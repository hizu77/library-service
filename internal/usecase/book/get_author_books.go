package book

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) GetAuthorBooks(ctx context.Context, id string) ([]entity.Book, error) {
	var outBooks []entity.Book

	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		var txErr error

		if _, txErr = u.authorRepository.GetAuthor(ctx, id); txErr != nil {
			u.logger.Error("authorRepository.GetAuthor", zap.String("id", id))
			return txErr
		}

		outBooks, txErr = u.bookRepository.GetBooksByAuthorID(ctx, id)
		if txErr != nil {
			u.logger.Error("bookRepository.GetBooksByAuthor", zap.String("id", id), zap.Error(txErr))
			return txErr
		}

		u.logger.Info("GetAuthorBooks", zap.String("ID", id))

		return nil
	})
	if err != nil {
		return nil, err
	}

	return outBooks, nil
}
