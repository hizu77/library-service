package httpserver

import (
	"net"
	"time"

	"go.uber.org/zap"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

func Prefork(prefork bool) Option {
	return func(s *Server) {
		s.prefork = prefork
	}
}

func Logger(logger *zap.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
