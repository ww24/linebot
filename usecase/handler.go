package usecase

import (
	"context"

	"github.com/ww24/linebot/domain/model"
)

type EventHandler interface {
	Handle(context.Context, []*model.Event) error
	HandleSchedule(context.Context) error
	HandleReminder(context.Context, *model.ReminderItemIDJSON) error
}
