package service

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
)

type Shopping interface {
	List(ctx context.Context, conversationID model.ConversationID) (model.ShoppingItems, error)
	AddItem(ctx context.Context, conversationID model.ConversationID, items ...*model.ShoppingItem) error
	DeleteAllItem(ctx context.Context, conversationID model.ConversationID) error
	DeleteItems(ctx context.Context, conversationID model.ConversationID, ids []string) error
	SetStatus(ctx context.Context, conversationID model.ConversationID) error
}

type ShoppingImpl struct {
	conversation repository.Conversation
	shopping     repository.Shopping
	tracer       trace.Tracer
}

func NewShopping(
	conversation repository.Conversation,
	shopping repository.Shopping,
	tracerProvider trace.TracerProvider,
) *ShoppingImpl {
	return &ShoppingImpl{
		conversation: conversation,
		shopping:     shopping,
		tracer:       tracerProvider.Tracer("github.com/ww24/linebot/domain/service"),
	}
}

func (s *ShoppingImpl) List(ctx context.Context, conversationID model.ConversationID) (model.ShoppingItems, error) {
	ctx, span := s.tracer.Start(ctx, "Shopping#List")
	defer span.End()

	if err := s.SetStatus(ctx, conversationID); err != nil {
		return nil, err
	}

	items, err := s.shopping.Find(ctx, conversationID)
	if err != nil {
		return nil, xerrors.Errorf("failed to find shopping items: %w", err)
	}

	return items, nil
}

func (s *ShoppingImpl) AddItem(ctx context.Context, conversationID model.ConversationID, items ...*model.ShoppingItem) error {
	ctx, span := s.tracer.Start(ctx, "Shopping#AddItem")
	defer span.End()

	if err := s.SetStatus(ctx, conversationID); err != nil {
		return err
	}
	if err := s.shopping.Add(ctx, items...); err != nil {
		return xerrors.Errorf("failed to add shopping item: %w", err)
	}
	return nil
}

func (s *ShoppingImpl) DeleteAllItem(ctx context.Context, conversationID model.ConversationID) error {
	ctx, span := s.tracer.Start(ctx, "Shopping#DeleteAllItem")
	defer span.End()

	if err := s.shopping.DeleteAll(ctx, conversationID); err != nil {
		return xerrors.Errorf("failed to delete all shopping items: %w", err)
	}
	if err := s.SetStatus(ctx, conversationID); err != nil {
		return err
	}
	return nil
}

func (s *ShoppingImpl) DeleteItems(ctx context.Context, conversationID model.ConversationID, ids []string) error {
	ctx, span := s.tracer.Start(ctx, "Shopping#DeleteItems")
	defer span.End()

	if err := s.shopping.BatchDelete(ctx, conversationID, ids); err != nil {
		return xerrors.Errorf("failed to delete shopping item: %w", err)
	}
	return nil
}

func (s *ShoppingImpl) SetStatus(ctx context.Context, conversationID model.ConversationID) error {
	ctx, span := s.tracer.Start(ctx, "Shopping#SetStatus")
	defer span.End()

	status := &model.ConversationStatus{
		ConversationID: conversationID,
		Type:           model.ConversationStatusTypeShopping,
	}
	if err := s.conversation.SetStatus(ctx, status); err != nil {
		return xerrors.Errorf("failed to set status: %w", err)
	}
	return nil
}
