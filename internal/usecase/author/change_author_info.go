package author

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) ChangeAuthorInfo(ctx context.Context, author entity.Author) (entity.Author, error) {
	span := trace.SpanFromContext(ctx)
	u.logger.Info(
		"Start to change author info",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("author_id", author.ID),
	)

	var outAuthor entity.Author
	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		var txErr error

		_, txErr = u.authorRepository.GetAuthor(ctx, author.ID)
		if txErr != nil {
			return txErr
		}

		outAuthor, txErr = u.authorRepository.UpdateAuthor(ctx, author)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		span.RecordError(err)

		u.logger.Error(
			"Change author info error",
			zap.Error(err),
			zap.String("author_id", outAuthor.ID),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)

		return entity.Author{}, err
	}

	span.SetAttributes(attribute.String("author_id", outAuthor.ID))
	u.logger.Info(
		"Change author info completed",
		zap.String("author_id", outAuthor.ID),
		zap.String("trace_id", span.SpanContext().TraceID().String()),
	)

	return outAuthor, nil
}
