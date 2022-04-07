package usecase

import (
	"context"
	"io"
	"net/url"

	"github.com/ww24/linebot/domain/model"
)

type EventHandler interface {
	Handle(context.Context, []*model.Event) error
	HandleSchedule(context.Context) error
	HandleReminder(context.Context, *model.ReminderItemIDJSON) error
}

type ScreenshotHandler interface {
	Handle(context.Context, *url.URL, string) (io.Reader, int, error)
}
