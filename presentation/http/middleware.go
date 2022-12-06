package http

import (
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
	"github.com/ww24/linebot/logger"
)

const trustedProxies = 1

func panicHandler(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("paniced in http handler", zap.Any("error", err))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func accessLogHandler(publisher accesslog.Publisher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			start := time.Now()

			accessLog := &avro.AccessLog{
				Timestamp:    &avro.UnionNullLong{Long: start.UnixMicro(), UnionType: avro.UnionNullLongTypeEnumLong},
				Id:           &avro.UnionNullString{String: xid.New().String(), UnionType: avro.UnionNullStringTypeEnumString},
				TraceId:      &avro.UnionNullString{},
				Ip:           &avro.UnionNullString{},
				UserAgent:    &avro.UnionNullString{String: r.UserAgent(), UnionType: avro.UnionNullStringTypeEnumString},
				Method:       &avro.UnionNullString{String: r.Method, UnionType: avro.UnionNullStringTypeEnumString},
				Path:         &avro.UnionNullString{String: r.URL.Path, UnionType: avro.UnionNullStringTypeEnumString},
				Query:        &avro.UnionNullString{String: r.URL.RawQuery, UnionType: avro.UnionNullStringTypeEnumString},
				Status:       &avro.UnionNullInt{},
				Duration:     &avro.UnionNullInt{},
				RequestSize:  &avro.UnionNullInt{Int: int32(r.ContentLength), UnionType: avro.UnionNullIntTypeEnumInt},
				ResponseSize: &avro.UnionNullInt{},
			}

			if t := trace.SpanContextFromContext(ctx).TraceID(); t.IsValid() {
				accessLog.TraceId.String = t.String()
				accessLog.TraceId.UnionType = avro.UnionNullStringTypeEnumString
			}

			ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
			if len(ips) > trustedProxies {
				clientIP := textproto.TrimString(ips[len(ips)-trustedProxies-1])
				accessLog.Ip.String = clientIP
				accessLog.Ip.UnionType = avro.UnionNullStringTypeEnumString
			}

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			accessLog.Status.Int = int32(rw.status)
			accessLog.Status.UnionType = avro.UnionNullIntTypeEnumInt
			accessLog.Duration.Int = int32(time.Since(start).Microseconds())
			accessLog.Duration.UnionType = avro.UnionNullIntTypeEnumInt
			accessLog.ResponseSize.Int = int32(rw.size)
			accessLog.ResponseSize.UnionType = avro.UnionNullIntTypeEnumInt

			publisher.Publish(ctx, accessLog)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func newResponseWriter(parent http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: parent}
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
