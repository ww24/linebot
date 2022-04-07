//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"net/http"

	"github.com/google/wire"

	"github.com/ww24/linebot/config"
	"github.com/ww24/linebot/infra/browser"
	"github.com/ww24/linebot/interactor"
	h "github.com/ww24/linebot/presentation/http"
)

func register(
	ctx context.Context,
) (*http.Server, error) {
	wire.Build(
		newLogger,
		config.Set,
		browser.Set,
		interactor.Set,
		h.ScreenshotHandlerSet,
		newServer,
	)
	return nil, nil
}
