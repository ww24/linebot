package tracer

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/log"
)

const shutdownTimeout = 5 * time.Second

func init() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// OpenCensus Bridge
	opencensus.InstallTraceBridge()
}

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

func New(c *Config, conf *config.Otel, exporter sdktrace.SpanExporter) (trace.TracerProvider, func()) {
	otel.SetLogger(
		logr.FromSlogHandler(
			log.NewLevelHandler(log.LevelFromEnv("OTEL_LOG_SEVERITY_LEVEL"), slog.Default().Handler()),
		),
	)

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(c.name),
		semconv.ServiceVersionKey.String(c.version),
	)

	sampler := newCustomSampler(
		sdktrace.ParentBased(sdktrace.TraceIDRatioBased(conf.SamplingRate),
			sdktrace.WithLocalParentSampled(sdktrace.AlwaysSample()),
			sdktrace.WithLocalParentNotSampled(sdktrace.NeverSample()),
			sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
			sdktrace.WithRemoteParentNotSampled(sdktrace.NeverSample()),
		),
		[]string{
			"google.devtools.cloudtrace.",
			"google.devtools.cloudprofiler.",
		},
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			slog.ErrorContext(ctx, "tracer: failed to shutdown Cloud Trace exporter", log.Err(err))
		}
	}
	return tp, cleanup
}

// customSampler implements sdktrace.Sampler.
var _ sdktrace.Sampler = (*customSampler)(nil)

type customSampler struct {
	parent                 sdktrace.Sampler
	ignoreSpanNamePrefixes []string
}

func newCustomSampler(parent sdktrace.Sampler, ignores []string) *customSampler {
	return &customSampler{parent: parent, ignoreSpanNamePrefixes: ignores}
}

//nolint:gocritic
func (s *customSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	for _, ignorePrefix := range s.ignoreSpanNamePrefixes {
		if strings.HasPrefix(p.Name, ignorePrefix) {
			return sdktrace.SamplingResult{
				Decision:   sdktrace.Drop,
				Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
			}
		}
	}
	return s.parent.ShouldSample(p)
}

func (s *customSampler) Description() string {
	return "CustomSampler"
}
