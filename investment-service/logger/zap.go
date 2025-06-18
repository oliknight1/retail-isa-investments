package logger

import (
	"os"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: l}
}

func NewLogger() (Logger, error) {
	var l *zap.Logger
	var err error

	switch os.Getenv("LOG_FORMAT") {
	case "json":
		l, err = zap.NewProduction()
	default:
		l, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return NewZapLogger(l), nil
}

func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, convert(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...Field) {
	z.logger.Error(msg, convert(fields)...)
}

func (z *ZapLogger) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, convert(fields)...)
}

func (z *ZapLogger) With(fields ...Field) Logger {
	return &ZapLogger{logger: z.logger.With(convert(fields)...)}
}

func (z *ZapLogger) Sync() error {
	return z.Sync()
}

// cast generic fields to zap fields
func convert(fields []Field) []zap.Field {
	zFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zFields[i] = f.(zap.Field)
	}
	return zFields
}
