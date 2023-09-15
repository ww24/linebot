package logger

import (
	"context"
	"sync/atomic"

	"github.com/blendle/zapdriver"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/internal/gcp"
)

//nolint:gochecknoglobals
var defaultLogger atomic.Pointer[Logger]

func init() {
	defaultLogger.Store(NewNop())
}

func SetConfig(name, version string) error {
	logger, err := new(name, version)
	if err != nil {
		return err
	}
	defaultLogger.Store(logger)
	return nil
}

func Default(ctx context.Context) *Logger {
	return defaultLogger.Load().WithTraceFromContext(ctx)
}

type Logger struct {
	*zap.Logger
	projectID string
}

func NewNop() *Logger { return &Logger{Logger: zap.NewNop()} }

func new(name, version string) (*Logger, error) {
	opt := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return newCore(core)
	})

	logger, err := zapdriver.NewProduction(opt)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize zapdriver: %w", err)
	}

	logger = logger.With(newServiceContext(name, version))

	projectID, err := gcp.ProjectID()
	if err != nil {
		logger.Warn("failed to get project id", zap.Error(err))
	}

	return &Logger{
		Logger:    logger,
		projectID: projectID,
	}, nil
}

func (l *Logger) WithTraceFromContext(ctx context.Context) *zap.Logger {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() || l.projectID == "" {
		return l.Logger
	}

	fields := zapdriver.TraceContext(
		spanCtx.TraceID().String(),
		spanCtx.SpanID().String(),
		spanCtx.IsSampled(),
		l.projectID,
	)
	return l.With(fields...)
}
