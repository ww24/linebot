package log

import (
	"context"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

func NewCloudLogging(w io.Writer, opts ...Option) *slog.Logger {
	o := &options{}
	o.apply(opts)

	logger := slog.New(newCloudLoggingHandler(w, o.gcpProjectID))
	if o.service != "" {
		logger = logger.With(serviceContextField(o.service, o.version))
	}
	if o.repository != "" && o.revisionID != "" {
		logger = logger.With(sourceReferenceField(o.repository, o.revisionID))
	}
	return logger
}

type CloudLoggingHandler struct {
	baseHandler slog.Handler
	projectID   string
}

func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.baseHandler.Enabled(ctx, level)
}

func (h *CloudLoggingHandler) Handle(ctx context.Context, r slog.Record) error {
	r = r.Clone()

	// add source location
	source := Source(r)
	r.AddAttrs(sourceLocationField(source))
	if r.Level >= slog.LevelError {
		r.AddAttrs(
			errorContextReportLocationField(source),
		)
	}
	r.AddAttrs(slog.String("stack_trace", StackTrace(r)))

	// add trace context
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() && h.projectID != "" {
		r.AddAttrs(
			traceContextFields(
				spanCtx.TraceID().String(),
				spanCtx.SpanID().String(),
				spanCtx.IsSampled(),
				h.projectID,
			)...,
		)
	}

	//nolint: wrapcheck
	return h.baseHandler.Handle(ctx, r)
}

func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CloudLoggingHandler{baseHandler: h.baseHandler.WithAttrs(attrs), projectID: h.projectID}
}

func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler {
	return &CloudLoggingHandler{baseHandler: h.baseHandler.WithGroup(name), projectID: h.projectID}
}

func newCloudLoggingHandler(w io.Writer, projectID GCPProjectID) slog.Handler {
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// top-level attributes
			if len(groups) == 0 {
				switch a.Key {
				case slog.TimeKey:
					a.Key = "timestamp"
				case slog.LevelKey:
					a.Key = "severity"
					if level, ok := a.Value.Any().(slog.Level); ok {
						a.Value = slog.StringValue(severity(level))
					}
				case slog.MessageKey:
					a.Key = "message"
				}
			}
			return a
		},
	})
	return &CloudLoggingHandler{baseHandler: h, projectID: string(projectID)}
}

func severity(level slog.Level) string {
	switch {
	case level < slog.LevelInfo:
		return "DEBUG"
	case level == slog.LevelInfo:
		return "INFO"
	case level < slog.LevelWarn:
		return "NOTICE"
	case level < slog.LevelError:
		return "WARNING"
	case level == slog.LevelError:
		return "ERROR"
	case level > slog.LevelError:
		return "CRITICAL"
	default:
		return "DEFAULT"
	}
}

// see: https://cloud.google.com/error-reporting/reference/rest/v1beta1/ServiceContext
// see: https://cloud.google.com/error-reporting/docs/formatting-error-messages
func serviceContextField(service Service, version Version) slog.Attr {
	fields := []any{
		slog.String("service", string(service)),
	}
	if version != "" {
		fields = append(fields, slog.String("version", string(version)))
	}
	return slog.Group("serviceContext", fields...)
}

// see: https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#sourcelocation
func errorContextReportLocationField(source *slog.Source) slog.Attr {
	return slog.Group("reportLocation",
		slog.String("filePath", source.File),
		slog.Int("lineNumber", source.Line),
		slog.String("functionName", source.Function),
	)
}

// see: https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext#SourceReference
func sourceReferenceField(repository Repository, revisionID RevisionID) slog.Attr {
	return slog.Any("sourceReference", []sourceReference{
		{
			Repository: strings.Replace(string(repository), "git://", "https://", 1),
			RevisionID: string(revisionID),
		},
	})
}

type sourceReference struct {
	Repository string `json:"repository"`
	RevisionID string `json:"revisionId"`
}

// see: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntrySourceLocation
func sourceLocationField(source *slog.Source) slog.Attr {
	return slog.Group("logging.googleapis.com/sourceLocation",
		slog.String("file", source.File),
		slog.String("line", strconv.Itoa(source.Line)),
		slog.String("function", source.Function),
	)
}

// see: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func traceContextFields(traceID, spanID string, sampled bool, project string) []slog.Attr {
	return []slog.Attr{
		slog.String("logging.googleapis.com/trace", "projects/"+project+"/traces/"+traceID),
		slog.String("logging.googleapis.com/spanId", spanID),
		slog.Bool("logging.googleapis.com/trace_sampled", sampled),
	}
}
