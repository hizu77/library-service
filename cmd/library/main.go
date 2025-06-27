package main

import (
	"github.com/hizu77/library-service/config"
	"github.com/hizu77/library-service/internal/app"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config error", zap.Error(err))
	}

	app.Run(cfg)
}
