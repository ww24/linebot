package tracer

import (
	"testing"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func TestMain(m *testing.M) {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagator.CloudTraceOneWayPropagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	m.Run()
}
