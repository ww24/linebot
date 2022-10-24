package tracer

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp/filters"
)

func HTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "server",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
			otelhttp.WithFilter(
				// ignore health check request
				filters.Not(filters.Path("/")),
			),
		)
	}
}

func HTTPTransport(parent http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(parent)
}
