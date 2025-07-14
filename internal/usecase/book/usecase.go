package book

import (
	"github.com/hizu77/library-service/internal/repository"
	"github.com/hizu77/library-service/internal/usecase"
	"go.uber.org/zap"
)

var _ usecase.BookUseCase = (*UseCaseImpl)(nil)

type UseCaseImpl struct {
	logger           *zap.Logger
	bookRepository   repository.BookRepository
	authorRepository repository.AuthorRepository
}

func NewUseCase(
	zap *zap.Logger,
	authorRepository repository.AuthorRepository,
	bookRepository repository.BookRepository) *UseCaseImpl {
	return &UseCaseImpl{
		logger:           zap,
		bookRepository:   bookRepository,
		authorRepository: authorRepository,
	}
}
