package postgres

import (
	"go.uber.org/zap"
	"time"
)

type Option func(*Postgres)

func MaxPoolSize(max int) Option {
	return func(p *Postgres) {
		p.maxPoolSize = max
	}
}

func MaxConnAttempts(max int) Option {
	return func(p *Postgres) {
		p.connAttempts = max
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
