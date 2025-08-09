package book

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hizu77/library-service/internal/infra/model/outbox"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) AddBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	span := trace.SpanFromContext(ctx)
	u.logger.Info(
		"Start to add book",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
	)

	var outBook entity.Book
	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		book.ID = uuid.New().String()
		var txErr error

		for i := range book.AuthorsIDs {
			_, txErr = u.authorRepository.GetAuthor(ctx, book.AuthorsIDs[i])
			if txErr != nil {
				return txErr
			}
		}

		outBook, txErr = u.bookRepository.AddBook(ctx, book)
		if txErr != nil {
			return txErr
		}

		serialized, txErr := json.Marshal(outBook)
		if txErr != nil {
			return txErr
		}

		idempotencyKey := "register_" + outbox.KindBook.String() + outBook.ID
		txErr = u.outboxRepository.SendMessage(ctx, idempotencyKey, outbox.KindBook, serialized)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		span.RecordError(err)

		u.logger.Error(
			"Add book error",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("book_id", outBook.ID),
		)

		return entity.Book{}, err
	}

	span.SetAttributes(attribute.String("book_id", outBook.ID))
	u.logger.Info(
		"Add book completed",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("book_id", outBook.ID),
	)

	return outBook, nil
}
