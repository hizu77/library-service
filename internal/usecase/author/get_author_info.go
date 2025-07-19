package author

import (
	"context"
	author "github.com/hizu77/library-service/internal/repository/persistent/author/postgres"

	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) GetAuthorInfo(ctx context.Context, id string) (entity.Author, error) {
	var outAuthor entity.Author

	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		var txErr error
		outAuthor, txErr = u.authorRepository.GetAuthor(ctx, id)
		if txErr != nil {
			u.logger.Error("authorRepository.GetAuthor", zap.Error(txErr))
			return txErr
		}

		u.logger.Info("GetAuthorInfo", zap.String("ID", author.ID))

		return nil
	})
	if err != nil {
		return entity.Author{}, err
	}

	return outAuthor, nil
}
