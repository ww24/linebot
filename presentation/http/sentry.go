package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"

	"github.com/ww24/linebot/internal/buildinfo"
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
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			const releasePrefix = "linebot@"
			if buildinfo.Version() != "" {
				event.Release = releasePrefix + buildinfo.Version()
			} else {
				event.Release = releasePrefix + "dev-" + buildinfo.Revision()
			}
			return event
		},
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

func report(r *http.Request, msg string, err error) {
	hub := sentry.GetHubFromContext(r.Context())
	if hub == nil {
		return
	}
	hub.WithScope(func(scope *sentry.Scope) {
		scope.SetExtra("message", msg)
		hub.CaptureException(err)
	})
}
