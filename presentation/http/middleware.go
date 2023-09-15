package http

import (
	"context"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/rs/xid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/internal/accesslog"
	"github.com/ww24/linebot/internal/accesslog/avro"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/logger"
)

func panicHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					ctx := r.Context()
					cl := logger.Default(ctx)
					cl.Error("paniced in http handler", zap.Any("error", err))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func accessLogHandler(publisher accesslog.Publisher, cfg *config.AccessLog) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			accessLog := &avro.AccessLog{
				Timestamp:    &avro.UnionNullLong{Long: start.UnixMicro(), UnionType: avro.UnionNullLongTypeEnumLong},
				Id:           &avro.UnionNullString{String: xid.New().String(), UnionType: avro.UnionNullStringTypeEnumString},
				TraceId:      nil,
				Ip:           nil,
				UserAgent:    &avro.UnionNullString{String: r.UserAgent(), UnionType: avro.UnionNullStringTypeEnumString},
				Method:       &avro.UnionNullString{String: r.Method, UnionType: avro.UnionNullStringTypeEnumString},
				Path:         &avro.UnionNullString{String: r.URL.Path, UnionType: avro.UnionNullStringTypeEnumString},
				Query:        &avro.UnionNullString{String: r.URL.RawQuery, UnionType: avro.UnionNullStringTypeEnumString},
				Status:       nil,
				Duration:     nil,
				RequestSize:  &avro.UnionNullInt{Int: int32(r.ContentLength), UnionType: avro.UnionNullIntTypeEnumInt},
				ResponseSize: nil,
			}

			if t := trace.SpanContextFromContext(r.Context()).TraceID(); t.IsValid() {
				accessLog.TraceId = &avro.UnionNullString{String: t.String(), UnionType: avro.UnionNullStringTypeEnumString}
			}

			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				ips := strings.Split(xff, ",")
				trustedProxies := cfg.TrustedProxies
				if len(ips) > trustedProxies {
					clientIP := textproto.TrimString(ips[len(ips)-trustedProxies-1])
					accessLog.Ip = &avro.UnionNullString{String: clientIP, UnionType: avro.UnionNullStringTypeEnumString}
				}
			}

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			accessLog.Status = &avro.UnionNullInt{Int: int32(rw.status), UnionType: avro.UnionNullIntTypeEnumInt}
			accessLog.Duration = &avro.UnionNullInt{Int: int32(time.Since(start).Microseconds()), UnionType: avro.UnionNullIntTypeEnumInt}
			accessLog.ResponseSize = &avro.UnionNullInt{Int: int32(rw.size), UnionType: avro.UnionNullIntTypeEnumInt}

			publisher.Publish(context.Background(), accessLog)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func newResponseWriter(parent http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: parent,
		status:         http.StatusOK,
		size:           0,
	}
}

func (w *responseWriter) Write(d []byte) (int, error) {
	n, err := w.ResponseWriter.Write(d)
	w.size += n
	if err != nil {
		return n, xerrors.Errorf("ResponseWriter.Write: %w", err)
	}
	return n, nil
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
