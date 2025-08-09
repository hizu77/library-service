package author

import (
	"context"

	author "github.com/hizu77/library-service/internal/repository/persistent/author/postgres"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) GetAuthorInfo(ctx context.Context, id string) (entity.Author, error) {
	span := trace.SpanFromContext(ctx)
	u.logger.Info(
		"Start to get author info",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("author_id", id),
	)

	var outAuthor entity.Author
	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		var txErr error
		outAuthor, txErr = u.authorRepository.GetAuthor(ctx, id)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		span.RecordError(err)

		u.logger.Error(
			"Get author info error",
			zap.Error(err),
			zap.String("author_id", outAuthor.ID),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)

		return entity.Author{}, err
	}

	span.SetAttributes(attribute.String("author_id", outAuthor.ID))
	u.logger.Info(
		"Get author info completed",
		zap.String("author_id", author.ID),
		zap.String("trace_id", span.SpanContext().TraceID().String()),
	)

	return outAuthor, nil
}
