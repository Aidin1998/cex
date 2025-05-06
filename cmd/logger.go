package cmd

import (
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(isDev bool) *zap.Logger {
	var zapLogger *zap.Logger

	if isDev {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapLogger = zap.Must(config.Build())
	} else {
		zapLogger = zap.Must(zap.NewProduction())
	}

	return zapLogger
}

func NewLogger(zapLogger *zap.Logger) *slog.Logger {
	return slog.New(zapslog.NewHandler(zapLogger.Core()))
}
