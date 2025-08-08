package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(logsPath string) (*zap.Logger, error) {
	err := os.MkdirAll(filepath.Dir(logsPath), 0755)
	if err != nil {
		return nil, errors.Wrap(err, "logger init")
	}

	file, err := os.OpenFile(logsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}

	writeSyncer := zapcore.AddSync(file)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)

	return zap.New(core), nil
}
