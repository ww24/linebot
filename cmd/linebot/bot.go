package main

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/tracer"
)

//nolint:gochecknoglobals
var tc = tracer.NewConfig(serviceName, version)

type bot struct {
	conf    *config.LINEBot
	handler http.Handler
}

func newBot(
	conf *config.LINEBot,
	handler http.Handler,
	_ trace.TracerProvider,
) *bot {
	return &bot{
		conf:    conf,
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
