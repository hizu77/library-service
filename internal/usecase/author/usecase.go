package author

import (
	"github.com/hizu77/library-service/internal/repository"
	"github.com/hizu77/library-service/internal/usecase"
	"go.uber.org/zap"
)

var _ usecase.AuthorUseCase = (*UseCaseImpl)(nil)

type UseCaseImpl struct {
	logger           *zap.Logger
	authorRepository repository.AuthorRepository
}

func NewUseCase(zap *zap.Logger, authorRepository repository.AuthorRepository) *UseCaseImpl {
	return &UseCaseImpl{
		logger:           zap,
		authorRepository: authorRepository,
	}
}
