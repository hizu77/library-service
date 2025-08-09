package author

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

func (u *UseCaseImpl) RegisterAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	span := trace.SpanFromContext(ctx)
	u.logger.Info(
		"Start to register author",
		zap.String("trace_id", span.SpanContext().TraceID().String()),
	)

	var outAuthor entity.Author
	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		author.ID = uuid.New().String()

		var txErr error
		outAuthor, txErr = u.authorRepository.AddAuthor(ctx, author)
		if txErr != nil {
			return txErr
		}

		serialized, txErr := json.Marshal(outAuthor)
		if txErr != nil {
			return txErr
		}

		idempotencyKey := "register_" + outbox.KindAuthor.String() + outAuthor.ID
		txErr = u.outboxRepository.SendMessage(ctx, idempotencyKey, outbox.KindAuthor, serialized)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		span.RecordError(err)

		u.logger.Error(
			"Register author error",
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("author_id", outAuthor.ID),
			zap.Error(err),
		)

		return entity.Author{}, err
	}

	span.SetAttributes(attribute.String("author_id", outAuthor.ID))
	u.logger.Info(
		"Register author completed",
		zap.String("author_id", outAuthor.ID),
		zap.String("trace_id", span.SpanContext().TraceID().String()),
	)

	return outAuthor, nil
}
