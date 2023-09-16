package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func spanContext(t *testing.T, traceIDHex, spanIDHex string) trace.SpanContext {
	t.Helper()
	traceID, err := trace.TraceIDFromHex(traceIDHex)
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex(spanIDHex)
	require.NoError(t, err)
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
}
