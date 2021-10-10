package service

import "github.com/google/wire"

// Set provides a wire set.
//nolint: gochecknoglobals
var Set = wire.NewSet(
	NewShopping,
	wire.Bind(new(Shopping), new(*ShoppingImpl)),
	NewBot,
	wire.Bind(new(Bot), new(*BotImpl)),
)
