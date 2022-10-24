package main

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/tracer"
)

//nolint:gochecknoglobals
var tc = tracer.NewConfig(serviceName, version)

type bot struct {
	config  repository.Config
	handler http.Handler
}

func newBot(
	config repository.Config,
	handler http.Handler,
	_ trace.TracerProvider,
) *bot {
	return &bot{
		config:  config,
		handler: handler,
	}
}

func newLogger(ctx context.Context) (*logger.Logger, error) {
	log, err := logger.New(ctx, serviceName, version)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize logger: %w", err)
	}

	return log, nil
}
