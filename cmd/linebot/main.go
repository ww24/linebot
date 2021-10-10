package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"cloud.google.com/go/profiler"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

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

	bot, err := register(ctx)
	if err != nil {
		log.Fatalf("initialize failed: %+v", err)
	}

	// set GOMAXPROCS
	if _, err := maxprocs.Set(maxprocs.Logger(bot.log.Sugar().Infof)); err != nil {
		bot.log.Warn("failed to set GOMAXPROCS", zap.Error(err))
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
			bot.log.Error("failed to start cloud profiler", zap.Error(err))
		}

		tp, err := tracer.New(serviceName, version)
		if err != nil {
			bot.log.Error("failed to initialize tracer", zap.Error(err))
		}
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				bot.log.Error("failed to shutdown tracer", zap.Error(err))
			}
		}()
	}

	addr := ":8000"
	if a := os.Getenv("PORT"); a != "" {
		addr = ":" + a
	}
	srv := &http.Server{
		Handler: bot.handler,
		Addr:    addr,
	}
	bot.log.Info("start server", zap.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)))
	//nolint:errcheck
	go srv.ListenAndServe()

	// wait signal
	<-ctx.Done()
	stop()

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(c); err != nil {
		bot.log.Error("failed to shutdown server", zap.Error(err))
	}
}
