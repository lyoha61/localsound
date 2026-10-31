package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.SugaredLogger, func(), error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zapLogger, err := config.Build()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = zapLogger.Sync()
	}

	return zapLogger.Sugar(), cleanup, nil
}
