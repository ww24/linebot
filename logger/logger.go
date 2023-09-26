package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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

func newLogger(w io.Writer, lvl zapcore.LevelEnabler) *Logger {
	ws := zapcore.Lock(zapcore.AddSync(w))
	enc := zapcore.NewJSONEncoder(newEncoderConfig())
	core := newCore(zapcore.NewCore(enc, ws, lvl))
	opts := []zap.Option{
		zap.AddCaller(),
	}
	logger := zap.New(core, opts...)

	projectID, err := gcp.ProjectID()
	if err != nil {
		logger.Warn("logger: failed to get project id", zap.Error(err))
	}

	return &Logger{
		Logger:    logger,
		projectID: projectID,
	}
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
	logger := newLogger(w, logLevel).withConfig(service, version)
	defaultLogger.Store(logger)
	return nil
}

func (l *Logger) WithTraceFromContext(ctx context.Context) *Logger {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() || l.projectID == "" {
		return l
	}

	fields := traceContextFields(
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

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func encodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var severity string
	switch l {
	case zapcore.DebugLevel:
		severity = "DEBUG"
	case zapcore.InfoLevel:
		severity = "INFO"
	case zapcore.WarnLevel:
		severity = "WARNING"
	case zapcore.ErrorLevel:
		severity = "ERROR"
	case zapcore.DPanicLevel:
		severity = "CRITICAL"
	case zapcore.PanicLevel:
		severity = "ALERT"
	case zapcore.FatalLevel:
		severity = "EMERGENCY"
	default:
		severity = "DEFAULT"
	}
	enc.AppendString(severity)
}

// see: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func traceContextFields(traceID, spanID string, sampled bool, project string) []zap.Field {
	return []zap.Field{
		zap.String("logging.googleapis.com/trace", fmt.Sprintf("projects/%s/traces/%s", project, traceID)),
		zap.String("logging.googleapis.com/spanId", spanID),
		zap.Bool("logging.googleapis.com/trace_sampled", sampled),
	}
}
