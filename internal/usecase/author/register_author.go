package author

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hizu77/library-service/internal/infra/model/outbox"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) RegisterAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	var outAuthor entity.Author

	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		author.ID = uuid.New().String()

		var txErr error
		outAuthor, txErr = u.authorRepository.AddAuthor(ctx, author)
		if txErr != nil {
			u.logger.Error("authorRepository.AddAuthor", zap.Error(txErr))
			return txErr
		}

		serialized, txErr := json.Marshal(outAuthor)
		if txErr != nil {
			u.logger.Error("Marshall", zap.Error(txErr))
			return txErr
		}

		idempotencyKey := "register_" + outbox.KindAuthor.String() + outAuthor.ID
		txErr = u.outboxRepository.SendMessage(ctx, idempotencyKey, outbox.KindAuthor, serialized)
		if txErr != nil {
			u.logger.Error("outboxRepository.SendMessage", zap.Error(txErr))
			return txErr
		}

		u.logger.Info("RegisterAuthor", zap.String("ID", outAuthor.ID))

		return nil
	})
	if err != nil {
		return entity.Author{}, err
	}

	return outAuthor, nil
}
