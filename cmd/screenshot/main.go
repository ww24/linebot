package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/profiler"
	"go.opentelemetry.io/otel"
	"go.uber.org/automaxprocs/maxprocs"
	"google.golang.org/grpc/grpclog"

	"github.com/ww24/linebot/internal/buildinfo"
	"github.com/ww24/linebot/internal/gcp"
	llog "github.com/ww24/linebot/log"
)

const (
	serviceName = "screenshot"
)

//nolint:gochecknoglobals
var tr = otel.Tracer("github.com/ww24/linebot/cmd/screenshot")

func init() {
	log.SetFlags(0)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, span := tr.Start(ctx, "start job")
	defer span.End()

	projectID, err := gcp.ProjectID()
	if err != nil {
		log.Printf("ERROR gcp.ProjectID: %+v", err)
	}

	if err := llog.SetOption(
		llog.Service(serviceName),
		llog.Version(buildinfo.Version()),
		llog.Repository(buildinfo.Repository()),
		llog.RevisionID(buildinfo.Revision()),
		llog.GCPProjectID(projectID),
	); err != nil {
		log.Printf("ERROR log.SetOption: %+v", err)
		stop()
		os.Exit(1)
	}
	grpclog.SetLoggerV2(llog.NewGRPCLogger(slog.Default().Handler()))

	// set GOMAXPROCS
	infof := func(format string, args ...interface{}) { slog.Info(fmt.Sprintf(format, args...)) }
	if _, err := maxprocs.Set(maxprocs.Logger(infof)); err != nil {
		slog.WarnContext(ctx, "main: failed to set GOMAXPROCS", llog.Err(err))
	}

	job, cleanup, err := register(ctx)
	if err != nil {
		stop()
		slog.ErrorContext(ctx, "main: register", llog.Err(err))
		os.Exit(1)
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
			slog.ErrorContext(ctx, "main: failed to start cloud profiler", llog.Err(err))
		}
	}

	slog.InfoContext(ctx, "main: start job")
	if err := job.run(ctx); err != nil {
		stop()
		cleanup()
		slog.ErrorContext(ctx, "main: failed to exec job", llog.Err(err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "main: done")
}
