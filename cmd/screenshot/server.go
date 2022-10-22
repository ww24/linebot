package main

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/tracer"
)

const readHeaderTimeout = 10 * time.Second

//nolint:gochecknoglobals
var tc = tracer.NewConfig(serviceName, version)

type server struct {
	config         repository.Config
	srv            *http.Server
	tracerProvider trace.TracerProvider
}

func newServer(
	conf repository.Config,
	handler http.Handler,
	tracerProvider trace.TracerProvider,
) *server {
	return &server{
		config: conf,
		srv: &http.Server{
			Handler:           handler,
			Addr:              conf.Addr(),
			ReadHeaderTimeout: readHeaderTimeout,
		},
		tracerProvider: tracerProvider,
	}
}

func newLogger(ctx context.Context) (*logger.Logger, error) {
	l, err := logger.New(ctx, serviceName, version)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize logger: %w", err)
	}

	return l, nil
}
