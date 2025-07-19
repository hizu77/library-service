package book

import (
	"github.com/hizu77/library-service/internal/repository"
	"github.com/hizu77/library-service/internal/usecase"
	"github.com/hizu77/library-service/pkg/transactor"
	"go.uber.org/zap"
)

var _ usecase.BookUseCase = (*UseCaseImpl)(nil)

type UseCaseImpl struct {
	logger           *zap.Logger
	bookRepository   repository.BookRepository
	authorRepository repository.AuthorRepository
	transactor       transactor.Transactor
}

func NewUseCase(
	zap *zap.Logger,
	authorRepository repository.AuthorRepository,
	bookRepository repository.BookRepository,
	transactor transactor.Transactor,
) *UseCaseImpl {
	return &UseCaseImpl{
		logger:           zap,
		bookRepository:   bookRepository,
		authorRepository: authorRepository,
		transactor:       transactor,
	}
}
