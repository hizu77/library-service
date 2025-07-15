package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

const (
	_defaultMaxPoolSize     = 1
	_defaultMaxConnAttempts = 10
	_defaultConnTimeout     = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
	logger  *zap.Logger
}

func New(ctx context.Context, url string, options ...Option) (*Postgres, error) {
	logger, _ := zap.NewDevelopment()

	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultMaxConnAttempts,
		connTimeout:  _defaultConnTimeout,
		logger:       logger,
	}

	for _, option := range options {
		option(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		logger.Error("pgxpool.ParseConfig", zap.Error(err))
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			break
		}

		logger.Error("pgxpool.NewWithConfig", zap.Error(err))

		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		logger.Error("attempts = 0", zap.Error(err))
		return nil, err
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
