package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/profiler"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/tracer"
)

const (
	serviceName     = "screenshot"
	shutdownTimeout = 10 * time.Second
)

var (
	// version is set during build
	version string
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.SetFlags(0)
	if err := logger.InitializeLogger(ctx, serviceName, version); err != nil {
		log.Printf("ERROR logger.InitializeLogger: %+v", err)
		return
	}
	dl := logger.DefaultLogger(ctx)

	// set GOMAXPROCS
	if _, err := maxprocs.Set(maxprocs.Logger(dl.Sugar().Infof)); err != nil {
		dl.Warn("failed to set GOMAXPROCS", zap.Error(err))
	}

	srv, err := register(ctx)
	if err != nil {
		dl.Error("register", zap.Error(err))
		panic(err)
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
			dl.Error("failed to start cloud profiler", zap.Error(err))
		}

		tp, err := tracer.New(serviceName, version, srv.config.OTELSamplingRate())
		if err != nil {
			dl.Error("failed to initialize tracer", zap.Error(err))
		}
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				dl.Error("failed to shutdown tracer", zap.Error(err))
			}
		}()
	}

	dl.Info("start server")
	//nolint:errcheck
	go srv.srv.ListenAndServe()

	// wait signal
	<-ctx.Done()
	stop()

	c, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.srv.Shutdown(c); err != nil {
		dl.Error("failed to shutdown server", zap.Error(err))
	}
}
