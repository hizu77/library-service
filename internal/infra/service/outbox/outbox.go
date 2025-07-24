package outbox

import (
	"context"
	"sync"
	"time"

	"github.com/hizu77/library-service/internal/infra/repository"
	"github.com/hizu77/library-service/internal/infra/service"
	"github.com/hizu77/library-service/pkg/transactor"
	"go.uber.org/zap"
)

var _ service.Outbox = (*Impl)(nil)

type Impl struct {
	outboxRepository repository.OutboxRepository
	logger           *zap.Logger
	transactor       transactor.Transactor
	GlobalHandler    GlobalHandler
	wg               *sync.WaitGroup
}

func (i *Impl) Stop() {
	i.wg.Wait()
	i.logger.Info("outbox shutting down")
}

func (i *Impl) Start(
	ctx context.Context,
	workers int,
	batchSize int,
	waitTime time.Duration,
	inProgressTTL time.Duration,
) {
	for worker := 0; worker < workers; worker++ {
		i.wg.Add(1)
		go i.worker(ctx, batchSize, waitTime, inProgressTTL)
	}
}

func (i *Impl) worker(
	ctx context.Context,
	batchSize int,
	waitTime time.Duration,
	inProgressTTL time.Duration,
) {
	ticker := time.NewTicker(waitTime)
	defer func() {
		ticker.Stop()
		i.wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := i.transactor.WithTx(ctx, func(ctx context.Context) error {
				messages, err := i.outboxRepository.GetMessages(ctx, batchSize, inProgressTTL)
				if err != nil {
					i.logger.Error("GetMessages", zap.Error(err))
					return err
				}

				successKeys := make([]string, 0, len(messages))
				for idx := range messages {
					message := messages[idx]

					kindHandler, err := i.GlobalHandler(message.Kind)
					if err != nil {
						i.logger.Error("GlobalHandler", zap.Error(err))
						continue
					}

					err = kindHandler(ctx, message.RawData)
					if err != nil {
						i.logger.Error("KindHandler", zap.Error(err))
						continue
					}

					successKeys = append(successKeys, message.IdempotencyKey)
				}

				err = i.outboxRepository.MarkMessageAsProcessed(ctx, successKeys)
				if err != nil {
					i.logger.Error("outboxRepository.MarkMessageAsProcessed", zap.Error(err))
					return err
				}

				return nil
			})

			if err != nil {
				i.logger.Error("Worker error", zap.Error(err))
			}
		}
	}
}

func New(
	outboxRepository repository.OutboxRepository,
	logger *zap.Logger,
	transactor transactor.Transactor,
	globalHandler GlobalHandler,
) *Impl {
	return &Impl{
		outboxRepository: outboxRepository,
		logger:           logger,
		transactor:       transactor,
		GlobalHandler:    globalHandler,
		wg:               &sync.WaitGroup{},
	}
}
