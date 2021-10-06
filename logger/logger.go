package logger

import (
	"context"
	"fmt"

	"github.com/blendle/zapdriver"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
)

type Logger struct {
	*zap.Logger
	projectID string
}

func New(ctx context.Context, name, version string) (*Logger, error) {
	core := zapdriver.WrapCore(
		zapdriver.ReportAllErrors(true),
		zapdriver.ServiceName(name),
	)

	logger, err := zapdriver.NewProductionWithCore(core)
	if err != nil {
		return nil, err
	}

	logger = logger.With(zap.String("version", version))

	projectID, err := getProjectID(ctx)
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

func getProjectID(ctx context.Context) (string, error) {
	cred, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to find default credentials: %w", err)
	}

	if cred.ProjectID == "" {
		return "", fmt.Errorf("project ID not found")
	}

	return cred.ProjectID, nil
}
