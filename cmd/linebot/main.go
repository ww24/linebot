package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"cloud.google.com/go/profiler"
	"google.golang.org/grpc/grpclog"

	"github.com/ww24/linebot/internal/buildinfo"
	"github.com/ww24/linebot/internal/gcp"
	llog "github.com/ww24/linebot/log"
)

const (
	serviceName       = "linebot"
	shutdownTimeout   = 10 * time.Second
	readHeaderTimeout = 10 * time.Second
)

func init() {
	log.SetFlags(0)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	projectID, err := gcp.ProjectID()
	if err != nil {
		panic(fmt.Errorf("ERROR gcp.ProjectID: %w", err))
	}

	if err := llog.SetOption(
		llog.Service(serviceName),
		llog.Version(buildinfo.Version()),
		llog.Repository(buildinfo.Repository()),
		llog.RevisionID(buildinfo.Revision()),
		llog.GCPProjectID(projectID),
	); err != nil {
		panic(fmt.Errorf("ERROR log.SetOption: %w", err))
	}
	grpclog.SetLoggerV2(llog.NewGRPCLogger(slog.Default().Handler()))

	bot, cleanup, err := register(ctx)
	if err != nil {
		slog.Error("main: register", llog.Err(err))
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
			slog.Error("main: failed to start cloud profiler", llog.Err(err))
		}
	}

	srv := &http.Server{
		Handler:           bot.handler,
		Addr:              bot.conf.Addr(),
		ReadHeaderTimeout: readHeaderTimeout,
	}
	slog.InfoContext(ctx, "main: start server",
		slog.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)),
		slog.String("addr", bot.conf.Addr()),
	)
	//nolint:errcheck
	go srv.ListenAndServe()

	// wait signal
	<-ctx.Done()
	stop()

	c, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(c); err != nil {
		slog.Error("main: failed to shutdown server", llog.Err(err))
	}
}
