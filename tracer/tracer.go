package tracer

import (
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"golang.org/x/xerrors"
)

func New(name, version string, samplingRate float64) (*sdktrace.TracerProvider, error) {
	exporter, err := cloudtrace.New()
	if err != nil {
		return nil, xerrors.Errorf("unable to set up tracing: %w", err)
	}

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(name),
		semconv.ServiceVersionKey.String(version),
	)

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplingRate),
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

	return tp, nil
}
