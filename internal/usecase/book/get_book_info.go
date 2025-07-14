package book

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) GetBookInfo(ctx context.Context, id string) (entity.Book, error) {
	v, err := u.bookRepository.GetBook(ctx, id)

	if err != nil {
		u.logger.Error("bookRepository.GetBook", zap.Error(err))
		return entity.Book{}, err
	}

	u.logger.Info("GetBookInfo", zap.String("ID", id))

	return v, nil
}
