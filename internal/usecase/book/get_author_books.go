package book

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) GetAuthorBooks(ctx context.Context, id string) ([]entity.Book, error) {
	span := trace.SpanFromContext(ctx)
	u.logger.Info(
		"Start to get author books",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("author_id", id),
	)

	var outBooks []entity.Book
	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		var txErr error

		if _, txErr = u.authorRepository.GetAuthor(ctx, id); txErr != nil {
			return txErr
		}

		outBooks, txErr = u.bookRepository.GetBooksByAuthorID(ctx, id)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		span.RecordError(err)

		u.logger.Error(
			"Get author books error",
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("author_id", id),
			zap.Error(err),
		)

		return nil, err
	}

	u.logger.Info(
		"Get author books completed",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("author_id", id),
	)

	return outBooks, nil
}
