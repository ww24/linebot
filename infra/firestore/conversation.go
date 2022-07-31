package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	gs "google.golang.org/grpc/status"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/internal/code"
)

type Conversation struct {
	*Client
}

func NewConversationStatus(src *model.ConversationStatus) *ConversationStatus {
	return &ConversationStatus{
		ConversationID: src.ConversationID,
		Status:         int(src.Type),
	}
}

func NewConversation(cli *Client) *Conversation {
	return &Conversation{Client: cli}
}

func (c *Conversation) conversations() *firestore.CollectionRef {
	return c.cli.Collection("conversations")
}

func (c *Conversation) conversation(conversationID model.ConversationID) *firestore.DocumentRef {
	return c.conversations().Doc(string(conversationID))
}

func (c *Conversation) SetStatus(ctx context.Context, status *model.ConversationStatus) error {
	ctx, span := tracer.Start(ctx, "Conversation#SetStatus")
	defer span.End()

	if err := status.Validate(); err != nil {
		return xerrors.Errorf("conversation status validation failed: %w", err)
	}

	conv := c.conversation(status.ConversationID)
	entity := NewConversationStatus(status)
	if _, err := conv.Set(ctx, entity); err != nil {
		return xerrors.Errorf("failed to set conversation status: %w", err)
	}

	return nil
}

func (c *Conversation) GetStatus(ctx context.Context, conversationID model.ConversationID) (*model.ConversationStatus, error) {
	ctx, span := tracer.Start(ctx, "Conversation#GetStatus")
	defer span.End()

	doc, err := c.conversation(conversationID).Get(ctx)
	if err != nil {
		if gs.Code(err) == codes.NotFound {
			return nil, code.With(err, code.NotFound)
		}
		return nil, xerrors.Errorf("failed to get conversation status: %w", err)
	}

	var ret ConversationStatus
	if err := doc.DataTo(&ret); err != nil {
		return nil, xerrors.Errorf("failed to convert response as ConversationStatus: %w", err)
	}
	return ret.Model(conversationID), nil
}

type ConversationStatus struct {
	ConversationID model.ConversationID `firestore:"-"`
	Status         int                  `firestore:"status"`
}

func (c *ConversationStatus) Model(conversationID model.ConversationID) *model.ConversationStatus {
	return &model.ConversationStatus{
		ConversationID: conversationID,
		Type:           model.ConversationStatusType(c.Status),
	}
}
