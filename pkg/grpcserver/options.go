package grpcserver

import (
	"net"

	"go.uber.org/zap"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

func Logger(logger *zap.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}
