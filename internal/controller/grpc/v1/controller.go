package v1

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/usecase"
	"go.uber.org/zap"
)

var _ generated.LibraryServer = (*ControllerImpl)(nil)

type ControllerImpl struct {
	zap           *zap.Logger
	bookUseCase   usecase.BookUseCase
	authorUseCase usecase.AuthorUseCase
}

func NewControllerImpl(
	logger *zap.Logger,
	bookUseCase usecase.BookUseCase,
	authorUseCase usecase.AuthorUseCase) *ControllerImpl {
	return &ControllerImpl{
		zap:           logger,
		bookUseCase:   bookUseCase,
		authorUseCase: authorUseCase,
	}
}
