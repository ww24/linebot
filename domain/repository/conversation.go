package repository

import (
	"context"

	"github.com/ww24/linebot/domain/model"
)

type Conversation interface {
	AddShoppingItem(ctx context.Context, item ...*model.ShoppingItem) error
	FindShoppingItem(ctx context.Context, conversationID string) ([]*model.ShoppingItem, error)
	DeleteShoppingItems(ctx context.Context, conversationID string, ids []string) error
	DeleteAllShoppingItem(ctx context.Context, conversationID string) error
	SetStatus(ctx context.Context, status *model.ConversationStatus) error
	GetStatus(ctx context.Context, conversationID string) (*model.ConversationStatus, error)
}
