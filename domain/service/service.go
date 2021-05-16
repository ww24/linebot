package service

import "github.com/google/wire"

var Set = wire.NewSet(
	NewShopping,
	wire.Bind(new(Shopping), new(*ShoppingImpl)),
)
