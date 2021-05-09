//+build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/ww24/linebot/bot"
	"github.com/ww24/linebot/infra/firestore"
)

func register(
	context.Context,
	bot.Config,
	firestore.ClientConfig,
) (*bot.Bot, error) {
	wire.Build(
		firestore.Set,
		bot.Set,
	)
	return nil, nil
}
