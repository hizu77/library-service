package app

import (
	"context"
	"net"
	"os/signal"
	"syscall"

	author "github.com/hizu77/library-service/internal/repository/persistent/author/inmemory"
	book "github.com/hizu77/library-service/internal/repository/persistent/book/inmemory"

	"github.com/hizu77/library-service/config"
	"github.com/hizu77/library-service/internal/controller/grpc"
	"github.com/hizu77/library-service/internal/controller/http"
	auc "github.com/hizu77/library-service/internal/usecase/author"
	buc "github.com/hizu77/library-service/internal/usecase/book"
	"github.com/hizu77/library-service/pkg/grpcserver"
	"github.com/hizu77/library-service/pkg/httpserver"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can not initialize logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	authorRepository := author.NewInMemoryRepository()
	bookRepository := book.NewInMemoryRepository()

	authorUseCase := auc.NewUseCase(logger, authorRepository)
	bookUseCase := buc.NewUseCase(logger, authorRepository, bookRepository)

	grpcServer := grpcserver.New(
		grpcserver.Port(cfg.GRPC.Port),
		grpcserver.Logger(logger))
	grpc.NewRouter(grpcServer.App, authorUseCase, bookUseCase, logger)

	gateway := httpserver.New(
		httpserver.Port(cfg.GRPC.GatewayPort),
		httpserver.Prefork(cfg.HTTP.UsePrefork))
	err = http.NewRouter(ctx, gateway.App, net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		logger.Fatal("can not initialize http server", zap.Error(err))
	}

	grpcServer.Start()
	gateway.Start()

	select {
	case <-ctx.Done():
		logger.Info("graceful shutdown")
	case err = <-gateway.Notify():
		logger.Error("Http server notify", zap.Error(err))
		cancel()
	case err = <-grpcServer.Notify():
		logger.Error("grpc server notify", zap.Error(err))
		cancel()
	}

	if err = gateway.Shutdown(); err != nil {
		logger.Error("gateway shutdown error", zap.Error(err))
	}

	if err = grpcServer.Shutdown(); err != nil {
		logger.Error("grpc server shutdown error", zap.Error(err))
	}
}
