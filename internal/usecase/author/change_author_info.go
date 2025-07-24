package author

import (
	"context"
	"encoding/json"
	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) ChangeAuthorInfo(ctx context.Context, author entity.Author) (entity.Author, error) {
	var outAuthor entity.Author

	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		var txErr error
		outAuthor, txErr = u.authorRepository.UpdateAuthor(ctx, author)
		if txErr != nil {
			u.logger.Error("authorRepository.UpdateAuthor", zap.Error(txErr))
			return txErr
		}

		serialized, txErr := json.Marshal(outAuthor)
		if txErr != nil {
			u.logger.Error("authorRepository.UpdateAuthor", zap.Error(txErr))
			return txErr
		}

		idempotencyKey :=

			u.logger.Info("ChangeAuthorInfo", zap.String("ID", outAuthor.ID))

		return nil

	})
	if err != nil {
		return entity.Author{}, err
	}

	return outAuthor, nil
}
