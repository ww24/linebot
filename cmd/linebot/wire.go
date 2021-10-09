//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/ww24/linebot/config"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/infra/external/linebot"
	"github.com/ww24/linebot/infra/firestore"
	"github.com/ww24/linebot/interactor"
	"github.com/ww24/linebot/nl"
	"github.com/ww24/linebot/presentation/http"
)

func register(
	context.Context,
) (*bot, error) {
	wire.Build(
		newLogger,
		config.Set,
		firestore.Set,
		linebot.Set,
		service.Set,
		nl.Set,
		interactor.Set,
		http.Set,
		newBot,
	)
	return nil, nil
}
