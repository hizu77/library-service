package repository

import (
	"context"
	"github.com/hizu77/library-service/internal/infra/model/outbox"
	"time"
)

type (
	OutboxRepository interface {
		SendMessage(ctx context.Context, idKey string, kind outbox.Kind, msg []byte) error
		GetMessages(ctx context.Context, batchSize int, inProgressTTL time.Duration) ([]outbox.Data, error)
		MarkMessageAsProcessed(ctx context.Context, idKeys []string) error
	}
)
