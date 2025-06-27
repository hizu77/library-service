package grpcserver

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	_defaultAddr = ":8080"
)

type Server struct {
	App     *grpc.Server
	notify  chan error
	address string
	logger  *zap.Logger
}

func New(opts ...Option) *Server {
	logger, _ := zap.NewDevelopment()

	server := &Server{
		App:     grpc.NewServer(),
		notify:  make(chan error),
		address: _defaultAddr,
		logger:  logger,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Start() {
	go func() {
		lis, err := net.Listen("tcp", s.address)

		if err != nil {
			s.logger.Error("grpc server failed to listen", zap.Error(err))
			s.notify <- err
			close(s.notify)

			return
		}

		s.logger.Info("grpc server listening at", zap.String("addr", s.address))

		if err = s.App.Serve(lis); err != nil {
			s.logger.Error("grpc server failed to serve", zap.Error(err))
			s.notify <- err
		}

		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	s.logger.Info("grpc server shutting down")
	s.App.GracefulStop()

	return nil
}
