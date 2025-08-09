package book

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) GetBookInfo(ctx context.Context, id string) (entity.Book, error) {
	span := trace.SpanFromContext(ctx)
	u.logger.Info(
		"Start to get book info",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("book_id", span.SpanContext().TraceID().String()),
	)

	var outBook entity.Book
	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		var txErr error
		outBook, txErr = u.bookRepository.GetBook(ctx, id)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		span.RecordError(err)

		u.logger.Error(
			"Get book info error",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("book_id", span.SpanContext().TraceID().String()),
		)

		return entity.Book{}, err
	}

	span.SetAttributes(attribute.String("book_id", outBook.ID))
	u.logger.Info(
		"Get book info completed",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("book_id", span.SpanContext().TraceID().String()),
	)

	return outBook, nil
}
