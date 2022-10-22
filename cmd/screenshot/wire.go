//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/ww24/linebot/config"
	"github.com/ww24/linebot/infra/browser"
	"github.com/ww24/linebot/interactor"
	h "github.com/ww24/linebot/presentation/http"
	"github.com/ww24/linebot/tracer"
)

func register(
	ctx context.Context,
) (*server, func(), error) {
	wire.Build(
		newLogger,
		config.Set,
		browser.Set,
		interactor.Set,
		h.ScreenshotHandlerSet,
		wire.Value(tc),
		tracer.NewCloudTraceExporter,
		tracer.New,
		newServer,
	)
	return nil, nil, nil
}
