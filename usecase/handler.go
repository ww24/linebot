package usecase

import (
	"context"
	"errors"

	"github.com/ww24/linebot/domain/model"
)

var ErrItemNotFound = errors.New("item not found")

type EventHandler interface {
	Handle(context.Context, []*model.Event) error
}
