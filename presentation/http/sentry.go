package http

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"

	"github.com/ww24/linebot/internal/config"
)

func newSentryMiddleware(cfg *config.Sentry) (*sentryhttp.Handler, error) {
	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           cfg.DSN,
		EnableTracing: false,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		return nil, fmt.Errorf("failed to initialize Sentry: %w", err)
	}

	// Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic:         true,
		WaitForDelivery: true,
		Timeout:         30 * time.Second,
	})

	return sentryHandler, nil
}
