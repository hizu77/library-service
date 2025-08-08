package httpserver

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

const (
	_defaultAddr            = ":8081"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	App    *fiber.App
	notify chan error
	logger *zap.Logger

	prefork         bool
	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
}

func New(logger *zap.Logger, opts ...Option) *Server {
	server := &Server{
		App:             nil,
		notify:          make(chan error, 1),
		logger:          logger,
		prefork:         false,
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(server)
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  server.readTimeout,
		WriteTimeout: server.writeTimeout,
		Prefork:      server.prefork,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	server.App = app

	return server
}

func (s *Server) Start() {
	go func() {
		s.logger.Info("http server listening on", zap.String("addr", s.address))

		if err := s.App.Listen(s.address); err != nil {
			s.logger.Error("http server failed to serve", zap.Error(err))
			s.notify <- err
		}

		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	s.logger.Info("http server shutting down")
	return s.App.ShutdownWithTimeout(s.shutdownTimeout)
}
