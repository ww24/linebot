package service

import (
	"context"
	"net/http"

	"golang.org/x/xerrors"

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
	events, err := b.bot.EventsFromRequest(r)
	if err != nil {
		return nil, xerrors.Errorf("failed to call EventsFromRequest: %w", err)
	}
	return events, nil
}

func (b *BotImpl) ReplyTextMessage(ctx context.Context, e *model.Event, text string) error {
	msg := b.message.Text(text)
	if err := b.bot.ReplyMessage(ctx, e, msg); err != nil {
		return xerrors.Errorf("failed to call ReplyMessage: %w", err)
	}
	return nil
}

func (b *BotImpl) ReplyMessage(ctx context.Context, e *model.Event, msg repository.MessageProvider) error {
	if err := b.bot.ReplyMessage(ctx, e, msg); err != nil {
		return xerrors.Errorf("failed to call ReplyMessage: %w", err)
	}
	return nil
}
