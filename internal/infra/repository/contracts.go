package repository

import (
	"context"
	"time"

	"github.com/hizu77/library-service/internal/infra/model/outbox"
)

//go:generate ../../../bin/mockgen -destination=../../usecase/mock/outbox.go -package=mock github.com/hizu77/library-service/internal/infra/repository OutboxRepository

type (
	OutboxRepository interface {
		SendMessage(ctx context.Context, idKey string, kind outbox.Kind, msg []byte) error
		GetMessages(ctx context.Context, batchSize int, inProgressTTL time.Duration) ([]outbox.Data, error)
		MarkMessageAsProcessed(ctx context.Context, idKeys []string) error
	}
)
