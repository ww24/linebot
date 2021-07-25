package repository

import (
	"context"
	"time"

	"github.com/ww24/linebot/domain/model"
)

type Executor interface {
	Do(context.Context, *model.ReminderItem, time.Time) error
}
