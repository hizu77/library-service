package book

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) UpdateBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	span := trace.SpanFromContext(ctx)
	u.logger.Info(
		"Start to update book",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("book_id", book.ID),
	)

	var outBook entity.Book
	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		var txErr error
		for i := range book.AuthorsIDs {
			_, txErr = u.authorRepository.GetAuthor(ctx, book.AuthorsIDs[i])
			if txErr != nil {
				return txErr
			}
		}

		outBook, txErr = u.bookRepository.UpdateBook(ctx, book)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		span.RecordError(err)

		u.logger.Error(
			"Update book error",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("book_id", book.ID),
		)

		return entity.Book{}, err
	}

	span.SetAttributes(attribute.String("book_id", outBook.ID))
	u.logger.Info(
		"Update book completed",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("book_id", span.SpanContext().TraceID().String()),
	)

	return outBook, nil
}
