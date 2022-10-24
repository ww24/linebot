package service

import (
	"github.com/google/wire"
	"go.opentelemetry.io/otel"
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewConversation,
	wire.Bind(new(Conversation), new(*ConversationImpl)),
	NewShopping,
	wire.Bind(new(Shopping), new(*ShoppingImpl)),
	NewReminder,
	wire.Bind(new(Reminder), new(*ReminderImpl)),
	NewBot,
	wire.Bind(new(Bot), new(*BotImpl)),
	NewWeather,
	wire.Bind(new(Weather), new(*WeatherImpl)),
)

var tracer = otel.Tracer("github.com/ww24/linebot/domain/service")
