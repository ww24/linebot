package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"

	"github.com/ww24/linebot/bot"
	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/tracer"
)

const (
	serviceName = "linebot"
)

var (
	// version is set during build
	version string
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	botCfg := bot.Config{
		ChannelSecret:   os.Getenv("LINE_CHANNEL_SECRET"),
		ChannelToken:    os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		ConversationIDs: strings.Split(os.Getenv("ALLOW_CONV_IDS"), ","),
	}
	bot, err := register(ctx, botCfg)
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	// initialize cloud profiler and tracing if build is production
	if version != "" {
		profilerConfig := profiler.Config{
			Service:           serviceName,
			ServiceVersion:    version,
			MutexProfiling:    true,
			EnableOCTelemetry: true,
		}
		if err := profiler.Start(profilerConfig); err != nil {
			// just log the error if failed to initialize profiler
			bot.Log.Error("failed to start cloud profiler", zap.Error(err))
		}

		tp, err := tracer.New(serviceName, version)
		if err != nil {
			bot.Log.Error("failed to initialize tracer", zap.Error(err))
		}
		tpShutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		defer func() {
			if err := tp.Shutdown(tpShutdownCtx); err != nil {
				bot.Log.Error("failed to shutdown tracer", zap.Error(err))
			}
		}()
	}

	addr := ":8000"
	if a := os.Getenv("PORT"); a != "" {
		addr = ":" + a
	}
	srv := &http.Server{
		Handler: handler(bot),
		Addr:    addr,
	}
	bot.Log.Info("start server")
	//nolint:errcheck
	go srv.ListenAndServe()

	// wait signal
	<-ctx.Done()
	stop()

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(c); err != nil {
		bot.Log.Error("failed to shutdown server", zap.Error(err))
	}
}

func newLogger(ctx context.Context) (*logger.Logger, error) {
	return logger.New(ctx, serviceName, version)
}

func handler(bot *bot.Bot) http.Handler {
	mux := http.NewServeMux()
	prop := propagator.New()
	mux.HandleFunc("/line_callback", func(w http.ResponseWriter, r *http.Request) {
		ctx := prop.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		cl := bot.Log.WithTraceFromContext(ctx)
		cl.Info("Request received")

		if err := bot.HandleRequest(ctx, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cl.Error("Request Error", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
	return mux
}
