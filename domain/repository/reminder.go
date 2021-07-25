package repository

import (
	"context"

	"github.com/ww24/linebot/domain/model"
)

type TaskScheduler interface {
	Add(context.Context, *model.ReminderItem) error
	List(context.Context, model.ConversationID) ([]string, error)
	Delete(context.Context, *model.ReminderItem) error
}
