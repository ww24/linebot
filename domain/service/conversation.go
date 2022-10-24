package service

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/internal/code"
)

type Conversation interface {
	GetStatus(context.Context, model.ConversationID) (*model.ConversationStatus, error)
	SetStatus(context.Context, *model.ConversationStatus) error
}

type ConversationImpl struct {
	conversation repository.Conversation
}

func NewConversation(
	conversation repository.Conversation,
) *ConversationImpl {
	return &ConversationImpl{
		conversation: conversation,
	}
}

func (s *ConversationImpl) GetStatus(ctx context.Context, conversationID model.ConversationID) (*model.ConversationStatus, error) {
	ctx, span := tracer.Start(ctx, "Conversation#GetStatus")
	defer span.End()

	status, err := s.conversation.GetStatus(ctx, conversationID)
	if code.From(err) == code.NotFound {
		status = &model.ConversationStatus{
			ConversationID: conversationID,
			Type:           model.ConversationStatusTypeNeutral,
		}
		err = s.conversation.SetStatus(ctx, status)
	}
	if err != nil {
		return nil, xerrors.Errorf("failed to get status: %w", err)
	}
	return status, nil
}

func (s *ConversationImpl) SetStatus(ctx context.Context, status *model.ConversationStatus) error {
	ctx, span := tracer.Start(ctx, "Conversation#SetStatus")
	defer span.End()

	if err := s.conversation.SetStatus(ctx, status); err != nil {
		return xerrors.Errorf("failed to set status: %w", err)
	}
	return nil
}
