package tracer

import (
	"context"
	"log/slog"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/ww24/linebot/log"
)

func NewCloudTraceExporter() sdktrace.SpanExporter {
	exporter, err := cloudtrace.New()
	if err != nil {
		slog.Error("tracer: unable to set up cloud trace exporter", log.Err(err))
		return new(noopExporter)
	}
	return exporter
}

// noopExporter implements sdktrace.SpanExporter
var _ sdktrace.SpanExporter = (*noopExporter)(nil)

type noopExporter struct{}

func (*noopExporter) ExportSpans(context.Context, []sdktrace.ReadOnlySpan) error { return nil }
func (*noopExporter) Shutdown(context.Context) error                             { return nil }
