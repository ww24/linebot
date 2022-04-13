//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOPACKAGE/mock_$GOFILE -package=mock_repository

package repository

import (
	"context"

	"github.com/ww24/linebot/domain/model"
)

type Conversation interface {
	SetStatus(context.Context, *model.ConversationStatus) error
	GetStatus(context.Context, model.ConversationID) (*model.ConversationStatus, error)
}
