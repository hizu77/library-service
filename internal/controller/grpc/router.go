package grpc

import (
	v1 "github.com/hizu77/library-service/internal/controller/grpc/v1"
	"github.com/hizu77/library-service/internal/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(
	server *grpc.Server,
	au usecase.AuthorUseCase,
	bu usecase.BookUseCase,
	logger *zap.Logger) {
	{
		v1.NewLibraryRoutes(server, au, bu, logger)
	}

	reflection.Register(server)
}
