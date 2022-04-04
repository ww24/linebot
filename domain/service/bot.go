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
	ReplyMessage(context.Context, *model.Event, repository.MessageProvider) error
	PushMessage(context.Context, model.ConversationID, repository.MessageProvider) error
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

func (b *BotImpl) ReplyMessage(ctx context.Context, e *model.Event, msg repository.MessageProvider) error {
	if err := b.bot.ReplyMessage(ctx, e, msg); err != nil {
		return xerrors.Errorf("failed to call ReplyMessage: %w", err)
	}
	return nil
}

func (b *BotImpl) PushMessage(ctx context.Context, conversationID model.ConversationID, msg repository.MessageProvider) error {
	if err := b.bot.PushMessage(ctx, conversationID, msg); err != nil {
		return xerrors.Errorf("failed to call PushMessage: %w", err)
	}
	return nil
}
