package service

import (
	"context"
	"net/http"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
)

type Bot interface {
	EventsFromRequest(r *http.Request) ([]*model.Event, error)
	ReplyTextMessage(context.Context, *model.Event, string) error
	ReplyMessage(context.Context, *model.Event, repository.MessageProvider) error
}

type BotImpl struct {
	bot     repository.Bot
	message repository.MessageProviderSet
}

func NewBot(
	bot repository.Bot,
	message repository.MessageProviderSet,
) *BotImpl {
	return &BotImpl{
		bot:     bot,
		message: message,
	}
}

func (b *BotImpl) EventsFromRequest(r *http.Request) ([]*model.Event, error) {
	return b.bot.EventsFromRequest(r)
}

func (b *BotImpl) ReplyTextMessage(ctx context.Context, e *model.Event, text string) error {
	msg := b.message.Text(text)
	return b.bot.ReplyMessage(ctx, e, msg)
}

func (b *BotImpl) ReplyMessage(ctx context.Context, e *model.Event, msg repository.MessageProvider) error {
	return b.bot.ReplyMessage(ctx, e, msg)
}
