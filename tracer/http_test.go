package tracer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestHTTPMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    *http.Request
		want trace.SpanContext
	}{
		{
			name: "request with Cloud Trace Context",
			r: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/trace", nil)
				r.Header.Set(propagator.TraceContextHeaderName, "105445aa7843bc8bf206b12000100000/1;o=1")
				return r
			}(),
			want: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    mustTraceIDFromHex("105445aa7843bc8bf206b12000100000"),
				SpanID:     mustSpanIDFromHex("0000000000000001"),
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			}),
		},
		{
			name: "ignore health check request",
			r: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.Header.Set(propagator.TraceContextHeaderName, "105445aa7843bc8bf206b12000100000/1;o=1")
				return r
			}(),
			want: trace.SpanContext{},
		},
		{
			name: "without Cloud Trace Context header",
			r: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/example", nil)
				return r
			}(),
			want: trace.SpanContext{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			HTTPMiddleware()(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ctx := r.Context()
					got := trace.SpanContextFromContext(ctx)
					assert.Equal(t, tt.want, got)
				}),
			).ServeHTTP(httptest.NewRecorder(), tt.r)
		})
	}
}

func mustTraceIDFromHex(h string) trace.TraceID {
	tid, err := trace.TraceIDFromHex(h)
	if err != nil {
		panic(err)
	}
	return tid
}

func mustSpanIDFromHex(h string) trace.SpanID {
	sid, err := trace.SpanIDFromHex(h)
	if err != nil {
		panic(err)
	}
	return sid
}
