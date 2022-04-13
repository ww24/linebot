//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOPACKAGE/mock_$GOFILE -package=mock_repository

package repository

import (
	"context"

	"github.com/ww24/linebot/domain/model"
)

type Shopping interface {
	Add(context.Context, ...*model.ShoppingItem) error
	Find(context.Context, model.ConversationID) ([]*model.ShoppingItem, error)
	BatchDelete(ctx context.Context, conversationID model.ConversationID, ids []string) error
	DeleteAll(context.Context, model.ConversationID) error
}
