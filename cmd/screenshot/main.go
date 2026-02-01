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
	if err := execute(context.Background()); err != nil {
		log.Printf("main: %+v", err)
		os.Exit(1)
	}
}

func execute(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, span := tr.Start(ctx, "start job")
	defer span.End()

	projectID, err := gcp.ProjectID()
	if err != nil {
		return fmt.Errorf("ERROR gcp.ProjectID: %w", err)
	}

	if err := llog.SetOption(
		llog.Service(serviceName),
		llog.Version(buildinfo.Version()),
		llog.Repository(buildinfo.Repository()),
		llog.RevisionID(buildinfo.Revision()),
		llog.GCPProjectID(projectID),
	); err != nil {
		return fmt.Errorf("ERROR log.SetOption: %w", err)
	}
	grpclog.SetLoggerV2(llog.NewGRPCLogger(slog.Default().Handler()))

	job, cleanup, err := register(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "main: register", llog.Err(err))
		return fmt.Errorf("ERROR register: %w", err)
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
		slog.ErrorContext(ctx, "main: failed to exec job", llog.Err(err))
		return fmt.Errorf("ERROR failed to exec job: %w", err)
	}

	slog.InfoContext(ctx, "main: done")
	return nil
}
