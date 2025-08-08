package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DirectoryPermissions = 0755
	FilePermissions      = 0644
)

func InitLogger(logsPath string) (*zap.Logger, error) {
	err := os.MkdirAll(
		filepath.Dir(logsPath),
		DirectoryPermissions,
	)
	if err != nil {
		return nil, errors.Wrap(err, "logger init")
	}

	file, err := os.OpenFile(
		logsPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		FilePermissions,
	)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}

	writeSyncer := zapcore.AddSync(file)
	consoleWroteSyncer := zapcore.AddSync(os.Stdout)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)
	consoleCore := zapcore.NewCore(encoder, consoleWroteSyncer, zap.InfoLevel)
	teeCore := zapcore.NewTee(core, consoleCore)

	return zap.New(teeCore), nil
}
