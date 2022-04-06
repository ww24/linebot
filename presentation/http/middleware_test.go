package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"

	"github.com/ww24/linebot/logger"
)

func TestPanicHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		h          http.Handler
		wantStatus int
		wantBody   string
	}{
		{
			name: "success handling",
			h: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "success")
			}),
			wantStatus: http.StatusOK,
			wantBody:   "success",
		},
		{
			name: "panic in handler",
			h: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("unexpected")
			}),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "Internal Server Error\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			PanicHandler(logger.NewNop())(tt.h).ServeHTTP(w, r)
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestXCTCOpenTelemetry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    *http.Request
		want trace.SpanContext
	}{
		{
			name: "request with Cloud Trace Context",
			r: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
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
			name: "general request",
			r: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				return r
			}(),
			want: trace.SpanContext{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			XCTCOpenTelemetry()(
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
