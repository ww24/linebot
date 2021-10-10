package usecase

import (
	"context"

	"github.com/ww24/linebot/domain/model"
)

type EventHandler interface {
	Handle(context.Context, []*model.Event) error
}
