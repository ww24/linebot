package repository

import (
	"context"
	"net/http"

	"github.com/ww24/linebot/domain/model"
)

type Bot interface {
	EventsFromRequest(r *http.Request) ([]*model.Event, error)
	ReplyMessage(context.Context, *model.Event, MessageProvider) error
}

type Handler interface {
	Handle(context.Context, *model.Event) error
}

type MessageProviderSet interface {
	Text(string) MessageProvider
	ShoppingDeleteConfirmation(string) MessageProvider
	ShoppingMenu(string, model.ShoppingReplyType) MessageProvider
}

type MessageProvider interface {
	AsMessage(interface{}) error
}
