package book

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hizu77/library-service/internal/infra/model/outbox"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) AddBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	var outBook entity.Book

	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		book.ID = uuid.New().String()
		var txErr error

		for i := range book.AuthorsIDs {
			_, txErr = u.authorRepository.GetAuthor(ctx, book.AuthorsIDs[i])
			if txErr != nil {
				u.logger.Error(
					"authorRepository.GetAuthor",
					zap.Error(txErr),
					zap.String("author_id", book.AuthorsIDs[i]),
				)

				return txErr
			}
		}

		outBook, txErr = u.bookRepository.AddBook(ctx, book)
		if txErr != nil {
			u.logger.Error(
				"bookRepository.AddBook",
				zap.Error(txErr),
				zap.String("book_id", book.ID),
			)

			return txErr
		}

		serialized, txErr := json.Marshal(outBook)
		if txErr != nil {
			u.logger.Error("Marshall", zap.Error(txErr))
			return txErr
		}

		idempotencyKey := "register_" + outbox.KindBook.String() + outBook.ID
		txErr = u.outboxRepository.SendMessage(ctx, idempotencyKey, outbox.KindBook, serialized)
		if txErr != nil {
			u.logger.Error("outboxRepository.SendMessage", zap.Error(txErr))
			return txErr
		}

		u.logger.Info(
			"AddBook",
			zap.String("book_id", outBook.ID),
		)

		return nil
	})
	if err != nil {
		return entity.Book{}, err
	}

	return outBook, nil
}
