package postgres

import (
	"time"

	"go.uber.org/zap"
)

type Option func(*Postgres)

func MaxPoolSize(maxSize int) Option {
	return func(p *Postgres) {
		p.maxPoolSize = maxSize
	}
}

func MaxConnAttempts(maxAttempts int) Option {
	return func(p *Postgres) {
		p.connAttempts = maxAttempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *Postgres) {
		p.connTimeout = timeout
	}
}

func Logger(logger *zap.Logger) Option {
	return func(s *Postgres) {
		s.logger = logger
	}
}
