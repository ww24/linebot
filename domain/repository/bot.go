//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOPACKAGE/mock_$GOFILE -package=mock_repository

package repository

import (
	"context"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/ww24/linebot/domain/model"
)

type Bot interface {
	EventsFromRequest(r *http.Request) ([]*model.Event, error)
	ReplyMessage(context.Context, *model.Event, MessageProvider) error
	PushMessage(context.Context, model.ConversationID, MessageProvider) error
}

type Handler interface {
	Handle(context.Context, *model.Event) error
}

type MessageProviderSet interface {
	Text(string) MessageProvider
	ShoppingDeleteConfirmation(string) MessageProvider
	ShoppingMenu(string, model.ShoppingReplyType) MessageProvider
	ReminderMenu(string, model.ReminderReplyType, []*model.ReminderItem) MessageProvider
	ReminderChoices(string, []string, []model.ExecutorType) MessageProvider
	TimePicker(text, data string) MessageProvider
	ReminderDeleteConfirmation(text, data string) MessageProvider
}

type MessageProvider interface {
	ToMessage() linebot.SendingMessage
}
