package tracer

import (
	"context"
	"time"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	octrace "go.opencensus.io/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/logger"
)

const shutdownTimeout = 5 * time.Second

type Config struct {
	name    string
	version string
}

func NewConfig(name, version string) *Config {
	return &Config{
		name:    name,
		version: version,
	}
}

func New(c *Config, conf repository.Config, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, func()) {
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(c.name),
		semconv.ServiceVersionKey.String(c.version),
	)

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(conf.OTELSamplingRate()),
		sdktrace.WithLocalParentSampled(sdktrace.AlwaysSample()),
		sdktrace.WithLocalParentNotSampled(sdktrace.NeverSample()),
		sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
		sdktrace.WithRemoteParentNotSampled(sdktrace.NeverSample()),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// OpenCensus Bridge
	tracer := otel.Tracer("OpenCensus")
	octrace.DefaultTracer = opencensus.NewTracer(tracer)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			dl := logger.DefaultLogger(ctx)
			dl.Error("failed to shutdown Cloud Trace exporter", zap.Error(err))
		}
	}
	return tp, cleanup
}

func NewCloudTraceExporter() (sdktrace.SpanExporter, error) {
	exporter, err := cloudtrace.New()
	if err != nil {
		return nil, xerrors.Errorf("unable to set up tracing: %w", err)
	}
	return exporter, nil
}
