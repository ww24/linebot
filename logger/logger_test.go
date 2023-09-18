package logger

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	defaultLogger.Load().projectID = "project-id"
	m.Run()
}

func TestLogger_withConfig(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		service string
		version string
		want    string
	}{
		{
			name: "no config",
			want: "",
		},
		{
			name:    "with service name and version",
			service: "service-name",
			version: "v1.0.0",
			want:    `"serviceContext":{"service":"service-name","version":"v1.0.0"}`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			l, err := newLogger(buf, zapcore.InfoLevel)
			require.NoError(t, err)
			if tt.service != "" {
				l = l.withConfig(tt.service, tt.version)
			}
			l.Info("withConfig")
			l.Sync()
			if tt.want == "" {
				assert.NotContains(t, buf.String(), `"serviceContext":`)
			} else {
				assert.Contains(t, buf.String(), tt.want)
			}
		})
	}
}

func TestLogger_WithTraceFromContext(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ctx  func(*testing.T) context.Context
		want string
	}{
		{
			name: "no trace",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
			want: "",
		},
		{
			name: "with trace",
			ctx: func(t *testing.T) context.Context {
				ctx := context.Background()
				const traceIDHex = "7e4ba55b36bb0d64c25dc7ac6d32a907"
				const spanIDHex = "05ce485f05506425"
				sc := spanContext(t, traceIDHex, spanIDHex)
				return trace.ContextWithSpanContext(ctx, sc)
			},
			want: `"logging.googleapis.com/trace":"projects/project-id/traces/7e4ba55b36bb0d64c25dc7ac6d32a907","logging.googleapis.com/spanId":"05ce485f05506425","logging.googleapis.com/trace_sampled":true,`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := &bytes.Buffer{}
			l, err := newLogger(buf, zapcore.InfoLevel)
			require.NoError(t, err)
			l.projectID = "project-id"
			l = l.WithTraceFromContext(tt.ctx(t))
			l.Info("WithTraceFromContext")
			l.Sync()
			if tt.want == "" {
				assert.NotContains(t, buf.String(), `"logging.googleapis.com/trace":`)
				assert.NotContains(t, buf.String(), `"logging.googleapis.com/spanId":`)
				assert.NotContains(t, buf.String(), `"logging.googleapis.com/trace_sampled":`)
			} else {
				assert.Contains(t, buf.String(), tt.want)
			}
		})
	}
}

func TestDefault_check_race(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ctx  func(*testing.T) context.Context
	}{
		{
			name: "no trace",
			ctx: func(t *testing.T) context.Context {
				return context.Background()
			},
		},
		{
			name: "with trace 1",
			ctx: func(t *testing.T) context.Context {
				ctx := context.Background()
				const traceIDHex = "7e4ba55b36bb0d64c25dc7ac6d32a907"
				const spanIDHex = "05ce485f05506425"
				sc := spanContext(t, traceIDHex, spanIDHex)
				return trace.ContextWithSpanContext(ctx, sc)
			},
		},
		{
			name: "with trace 2",
			ctx: func(t *testing.T) context.Context {
				ctx := context.Background()
				const traceIDHex = "242c372fb5b39bf4b518e880ed57db37"
				const spanIDHex = "920dbfadc59271e0"
				sc := spanContext(t, traceIDHex, spanIDHex)
				return trace.ContextWithSpanContext(ctx, sc)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			Default(tt.ctx(t))
		})
	}
}

func TestSetConfig_check_race(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		service string
		version string
	}{
		{
			name:    "no version",
			service: "service-name",
		},
		{
			name:    "service name and version",
			service: "service-name",
			version: "v1.0.0",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			SetConfig(tt.service, tt.version)
		})
	}
}
