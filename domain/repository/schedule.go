//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOPACKAGE/mock_$GOFILE -package=mock_repository

package repository

import (
	"context"
	"time"

	"github.com/ww24/linebot/domain/model"
)

type ScheduleHandler interface {
	HandleSchedule(context.Context) error
}

type ScheduleSynchronizer interface {
	Sync(context.Context, model.ConversationID, model.ReminderItems, time.Time) error
	Create(context.Context, model.ConversationID, *model.ReminderItem, time.Time) error
	Delete(context.Context, model.ConversationID, *model.ReminderItem, time.Time) error
}

type RemindHandler interface {
	HandleReminder(context.Context, *model.ReminderItem) error
}
