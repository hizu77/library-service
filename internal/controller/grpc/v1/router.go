package v1

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewLibraryRoutes(
	server *grpc.Server,
	au usecase.AuthorUseCase,
	bu usecase.BookUseCase,
	logger *zap.Logger) {
	ctrl := NewControllerImpl(logger, bu, au)

	{
		generated.RegisterLibraryServer(server, ctrl)
	}
}
