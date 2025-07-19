package author

import (
	"github.com/hizu77/library-service/internal/repository"
	"github.com/hizu77/library-service/internal/usecase"
	"github.com/hizu77/library-service/pkg/transactor"
	"go.uber.org/zap"
)

var _ usecase.AuthorUseCase = (*UseCaseImpl)(nil)

type UseCaseImpl struct {
	logger           *zap.Logger
	authorRepository repository.AuthorRepository
	transactor       transactor.Transactor
}

func NewUseCase(
	zap *zap.Logger,
	authorRepository repository.AuthorRepository,
	transactor transactor.Transactor,
) *UseCaseImpl {
	return &UseCaseImpl{
		logger:           zap,
		authorRepository: authorRepository,
		transactor:       transactor,
	}
}
