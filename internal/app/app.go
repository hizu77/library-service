package app

import (
	"context"
	"net"
	"os/signal"
	"syscall"

	"github.com/hizu77/library-service/db"
	"github.com/hizu77/library-service/internal/bootstrap"
	authorRepo "github.com/hizu77/library-service/internal/repository/persistent/author/postgres"
	"github.com/hizu77/library-service/pkg/postgres"
	"github.com/hizu77/library-service/pkg/transactor"

	outboxRepo "github.com/hizu77/library-service/internal/infra/repository/outbox"
	outboxService "github.com/hizu77/library-service/internal/infra/service/outbox"
	bookRepo "github.com/hizu77/library-service/internal/repository/persistent/book/postgres"

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
	logger, err := bootstrap.InitLogger(cfg.Logger.LogFilePath)
	if err != nil {
		log.Fatalf("can not initialize logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pg, err := postgres.New(ctx, cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		logger.Fatal("can not initialize postgres", zap.Error(err))
	}
	defer pg.Close()

	db.Migrate(pg.Pool, logger)

	authorRepository := authorRepo.New(pg)
	bookRepository := bookRepo.New(pg)
	outboxRepository := outboxRepo.New(pg)

	tx := transactor.New(pg.Pool)

	outbox := outboxService.New(outboxRepository, logger, tx, outboxService.Handler())
	authorUseCase := auc.NewUseCase(logger, authorRepository, tx, outboxRepository)
	bookUseCase := buc.NewUseCase(logger, authorRepository, bookRepository, tx, outboxRepository)

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
	outbox.Start(
		ctx,
		cfg.Outbox.Workers,
		cfg.Outbox.BatchSize,
		cfg.Outbox.WaitTimeMS,
		cfg.Outbox.InProgressTTLMS,
	)

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

	outbox.Stop()
}
