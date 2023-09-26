package tracer

import (
	"context"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/ww24/linebot/logger"
)

func NewCloudTraceExporter() sdktrace.SpanExporter {
	exporter, err := cloudtrace.New()
	if err != nil {
		dl := logger.Default(context.Background())
		dl.Error("tracer: unable to set up cloud trace exporter", zap.Error(err))
		return new(noopExporter)
	}
	return exporter
}

// noopExporter implements sdktrace.SpanExporter
var _ sdktrace.SpanExporter = (*noopExporter)(nil)

type noopExporter struct{}

func (*noopExporter) ExportSpans(context.Context, []sdktrace.ReadOnlySpan) error { return nil }
func (*noopExporter) Shutdown(context.Context) error                             { return nil }
