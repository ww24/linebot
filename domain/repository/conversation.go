package repository

import (
	"context"

	"github.com/ww24/linebot/domain/model"
)

type Conversation interface {
	AddShoppingItem(context.Context, ...*model.ShoppingItem) error
	FindShoppingItem(context.Context, model.ConversationID) ([]*model.ShoppingItem, error)
	DeleteShoppingItems(ctx context.Context, conversationID model.ConversationID, ids []string) error
	DeleteAllShoppingItem(context.Context, model.ConversationID) error
	SetStatus(context.Context, *model.ConversationStatus) error
	GetStatus(context.Context, model.ConversationID) (*model.ConversationStatus, error)
}

type Reminder interface {
	Add(context.Context, model.ConversationID, *model.ReminderItem) error
	Find(context.Context, model.ConversationID) ([]*model.ReminderItem, error)
	Delete(context.Context, model.ConversationID, string) error
}
