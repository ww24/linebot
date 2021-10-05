//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/ww24/linebot/bot"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/infra/firestore"
)

func register(
	context.Context,
	bot.Config,
) (*bot.Bot, error) {
	wire.Build(
		newLogger,
		firestore.Set,
		service.Set,
		bot.Set,
	)
	return nil, nil
}
