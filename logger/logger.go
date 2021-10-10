package logger

import (
	"context"

	"github.com/blendle/zapdriver"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2/google"
	"golang.org/x/xerrors"
)

type Logger struct {
	*zap.Logger
	projectID string
}

func New(ctx context.Context, name, version string) (*Logger, error) {
	core := zapdriver.WrapCore(
		zapdriver.ReportAllErrors(true),
	)

	logger, err := zapdriver.NewProductionWithCore(core)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize zapdriver: %w", err)
	}

	logger = logger.With(
		zap.Object("serviceContext", newServiceContext(name, version)),
	)

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
		return "", xerrors.Errorf("failed to find default credentials: %w", err)
	}

	if cred.ProjectID == "" {
		return "", xerrors.Errorf("project ID not found")
	}

	return cred.ProjectID, nil
}

// see: https://cloud.google.com/error-reporting/reference/rest/v1beta1/ServiceContext
// see: https://cloud.google.com/error-reporting/docs/formatting-error-messages
type serviceContext struct {
	service string
	version string
}

func newServiceContext(service, version string) *serviceContext {
	return &serviceContext{
		service: service,
		version: version,
	}
}

func (c *serviceContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("service", c.service)
	if c.version != "" {
		enc.AddString("version", c.version)
	}
	return nil
}
