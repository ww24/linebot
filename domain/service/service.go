package service

import "github.com/google/wire"

// Set provides a wire set.
var Set = wire.NewSet(
	NewConversation,
	wire.Bind(new(Conversation), new(*ConversationImpl)),
	NewShopping,
	wire.Bind(new(Shopping), new(*ShoppingImpl)),
	NewBot,
	wire.Bind(new(Bot), new(*BotImpl)),
)
