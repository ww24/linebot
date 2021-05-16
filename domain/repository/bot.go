package repository

import (
	"context"

	"github.com/ww24/linebot/domain/model"
)

type Bot interface {
	PostMessage(context.Context, model.ConversationID, string) error
}
