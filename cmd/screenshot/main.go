package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/profiler"
	"go.opentelemetry.io/otel"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

	"github.com/ww24/linebot/logger"
)

const (
	serviceName = "screenshot"
)

var (
	// version is set during build
	version string
	//nolint:gochecknoglobals
	tr = otel.Tracer("github.com/ww24/linebot/cmd/screenshot")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, span := tr.Start(ctx, "start job")
	defer span.End()

	log.SetFlags(0)
	if err := logger.SetConfig(serviceName, version); err != nil {
		stop()
		log.Printf("ERROR logger.SetMeta: %+v", err)
		os.Exit(1)
	}
	dl := logger.DefaultLogger(ctx)

	// set GOMAXPROCS
	if _, err := maxprocs.Set(maxprocs.Logger(dl.Sugar().Infof)); err != nil {
		dl.Warn("failed to set GOMAXPROCS", zap.Error(err))
	}

	job, cleanup, err := register(ctx)
	if err != nil {
		stop()
		dl.Error("register", zap.Error(err))
		os.Exit(1)
	}
	defer cleanup()

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
	}

	dl.Info("start job")
	if err := job.run(ctx); err != nil {
		stop()
		cleanup()
		dl.Error("failed to exec job", zap.Error(err))
		os.Exit(1)
	}

	dl.Info("done")
}
