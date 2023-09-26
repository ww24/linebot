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

	"github.com/ww24/linebot/internal/buildinfo"
	"github.com/ww24/linebot/logger"
)

const (
	serviceName       = "linebot"
	shutdownTimeout   = 10 * time.Second
	readHeaderTimeout = 10 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.SetFlags(0)
	if err := logger.SetConfig(serviceName, buildinfo.Version()); err != nil {
		log.Printf("ERROR logger.SetMeta: %+v", err)
		return
	}
	dl := logger.Default(ctx)

	// set GOMAXPROCS
	if _, err := maxprocs.Set(maxprocs.Logger(dl.Sugar().Infof)); err != nil {
		dl.Warn("main: failed to set GOMAXPROCS", zap.Error(err))
	}

	bot, cleanup, err := register(ctx)
	if err != nil {
		dl.Error("main: register", zap.Error(err))
		panic(err)
	}
	defer cleanup()

	// initialize cloud profiler and tracing if build is production
	if version := buildinfo.Version(); version != "" {
		profilerConfig := profiler.Config{
			Service:           serviceName,
			ServiceVersion:    version,
			MutexProfiling:    true,
			EnableOCTelemetry: true,
		}
		if err := profiler.Start(profilerConfig); err != nil {
			// just log the error if failed to initialize profiler
			dl.Error("main: failed to start cloud profiler", zap.Error(err))
		}
	}

	srv := &http.Server{
		Handler:           bot.handler,
		Addr:              bot.conf.Addr(),
		ReadHeaderTimeout: readHeaderTimeout,
	}
	dl.Info("main: start server",
		zap.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)),
		zap.String("addr", bot.conf.Addr()),
	)
	//nolint:errcheck
	go srv.ListenAndServe()

	// wait signal
	<-ctx.Done()
	stop()

	c, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(c); err != nil {
		dl.Error("main: failed to shutdown server", zap.Error(err))
	}
}
