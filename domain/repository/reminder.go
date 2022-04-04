//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOPACKAGE/mock_$GOFILE -package=mock_repository

package repository

import (
	"context"

	"github.com/ww24/linebot/domain/model"
)

type Reminder interface {
	Add(context.Context, *model.ReminderItem) error
	List(context.Context, model.ConversationID) ([]*model.ReminderItem, error)
	Get(context.Context, model.ConversationID, model.ReminderItemID) (*model.ReminderItem, error)
	Delete(context.Context, model.ConversationID, model.ReminderItemID) error
	ListAll(context.Context) ([]*model.ReminderItem, error)
}
