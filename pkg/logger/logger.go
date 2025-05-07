package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	*zap.Logger
}

func New(logLevel string) (*Logger, error) {
	lvl := zap.InfoLevel
	if err := lvl.Set(logLevel); err != nil {
		return nil, fmt.Errorf("setting log level: %w", err)
	}

	zapLogger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.Lock(os.Stderr),
			zap.NewAtomicLevelAt(lvl)),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.LevelOf(zap.ErrorLevel)),
	)

	// zapLogger, err := zap.New()
	return &Logger{
		zapLogger,
	}, nil
}

type ctxKey struct{}

var defaultContextLogger *Logger

func (l *Logger) WithContext(ctx context.Context) context.Context {
	if _, ok := ctx.Value(ctxKey{}).(*Logger); ok {
		return ctx
	}

	return context.WithValue(ctx, ctxKey{}, l)
}

func Ctx(ctx context.Context) *Logger {
	if l, ok := ctx.Value(ctxKey{}).(*Logger); ok {
		return l
	}

	return defaultContextLogger
}
