package logger

import (
	"context"
	"io"
	"os"
	"sync/atomic"

	"github.com/blendle/zapdriver"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/internal/gcp"
)

const logLevelEnv = "LOG_LEVEL"

//nolint:gochecknoglobals
var defaultLogger atomic.Pointer[Logger]

func init() {
	defaultLogger.Store(NewNop())
}

func Default(ctx context.Context) *Logger {
	return defaultLogger.Load().WithTraceFromContext(ctx)
}

type Logger struct {
	*zap.Logger
	projectID string
}

func NewNop() *Logger { return &Logger{Logger: zap.NewNop()} }

func newLogger(w io.Writer, lvl zapcore.LevelEnabler) (*Logger, error) {
	ws := zapcore.Lock(zapcore.AddSync(w))
	enc := zapcore.NewJSONEncoder(zapdriver.NewProductionEncoderConfig())
	opts := []zap.Option{
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			core := zapcore.NewCore(enc, ws, lvl)
			return newCore(core)
		}),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	logger, err := zapdriver.NewProduction(opts...)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize zapdriver: %w", err)
	}

	projectID, err := gcp.ProjectID()
	if err != nil {
		logger.Warn("failed to get project id", zap.Error(err))
	}

	return &Logger{
		Logger:    logger,
		projectID: projectID,
	}, nil
}

func (l *Logger) clone() *Logger {
	cp := *l
	return &cp
}

func (l *Logger) withConfig(service, version string) *Logger {
	return l.WithLogger(l.With(newServiceContext(service, version)))
}

func (l *Logger) WithLogger(zl *zap.Logger) *Logger {
	cp := l.clone()
	cp.Logger = zl
	return cp
}

func SetConfig(service, version string) error {
	return SetConfigWithWriter(service, version, os.Stderr)
}

func SetConfigWithWriter(service, version string, w io.Writer) error {
	logLevel := getLogLevel(logLevelEnv)
	logger, err := newLogger(w, logLevel)
	if err != nil {
		return err
	}
	logger = logger.withConfig(service, version)
	defaultLogger.Store(logger)
	return nil
}

func (l *Logger) WithTraceFromContext(ctx context.Context) *Logger {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() || l.projectID == "" {
		return l
	}

	fields := zapdriver.TraceContext(
		spanCtx.TraceID().String(),
		spanCtx.SpanID().String(),
		spanCtx.IsSampled(),
		l.projectID,
	)
	return l.WithLogger(l.With(fields...))
}

func getLogLevel(env string) zapcore.Level {
	logLevel := os.Getenv(env)
	switch logLevel {
	case "DEBUG", "debug":
		return zapcore.DebugLevel
	case "INFO", "info":
		return zapcore.InfoLevel
	case "WARN", "warn":
		return zapcore.WarnLevel
	case "ERROR", "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
