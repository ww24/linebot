package main

import (
	"context"
	"net/url"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/config"
	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/tracer"
	"github.com/ww24/linebot/usecase"
)

//nolint:gochecknoglobals
var tc = tracer.NewConfig(serviceName, version)

type job struct {
	config     *config.Screenshot
	screenshot usecase.ScreenshotHandler
}

func newJob(
	conf *config.Screenshot,
	screenshot usecase.ScreenshotHandler,
	_ trace.TracerProvider,
	_ *logger.Logger,
) *job {
	return &job{
		config:     conf,
		screenshot: screenshot,
	}
}

func (j *job) run(ctx context.Context) error {
	target, err := url.Parse(j.config.TargetURL)
	if err != nil {
		return xerrors.Errorf("failed to parse target url: %w", err)
	}

	if err := j.screenshot.Handle(ctx, target, j.config.TargetSelector); err != nil {
		return xerrors.Errorf("failed to handle screenshot: %w", err)
	}

	return nil
}

func newLogger(ctx context.Context) (*logger.Logger, error) {
	l, err := logger.New(ctx, serviceName, version)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize logger: %w", err)
	}

	return l, nil
}
