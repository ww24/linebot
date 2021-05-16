package service

import (
	"context"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"golang.org/x/xerrors"
)

type Shopping interface {
	List(ctx context.Context, conversationID model.ConversationID) (model.ShoppingItems, error)
	AddItem(ctx context.Context, conversationID model.ConversationID, items ...*model.ShoppingItem) error
	DeleteAllItem(ctx context.Context, conversationID model.ConversationID) error
	DeleteItems(ctx context.Context, conversationID model.ConversationID, ids []string) error
	GetStatus(ctx context.Context, conversationID model.ConversationID) (*model.ConversationStatus, error)
	SetStatus(ctx context.Context, status *model.ConversationStatus) error
	SetStatusShopping(ctx context.Context, conversationID model.ConversationID) error
}

type ShoppingImpl struct {
	conversation repository.Conversation
}

func NewShopping(conversation repository.Conversation) (*ShoppingImpl, error) {
	return &ShoppingImpl{
		conversation: conversation,
	}, nil
}

func (s *ShoppingImpl) List(ctx context.Context, conversationID model.ConversationID) (model.ShoppingItems, error) {
	if err := s.SetStatusShopping(ctx, conversationID); err != nil {
		return nil, err
	}

	items, err := s.conversation.FindShoppingItem(ctx, conversationID)
	if err != nil {
		return nil, xerrors.Errorf("failed to find shopping items: %w", err)
	}

	return items, nil
}

func (s *ShoppingImpl) AddItem(ctx context.Context, conversationID model.ConversationID, items ...*model.ShoppingItem) error {
	if err := s.SetStatusShopping(ctx, conversationID); err != nil {
		return err
	}
	if err := s.conversation.AddShoppingItem(ctx, items...); err != nil {
		return xerrors.Errorf("failed to add shopping item: %w", err)
	}
	return nil
}

func (s *ShoppingImpl) DeleteAllItem(ctx context.Context, conversationID model.ConversationID) error {
	if err := s.conversation.DeleteAllShoppingItem(ctx, conversationID); err != nil {
		return xerrors.Errorf("failed to delete all shopping items: %w", err)
	}
	if err := s.SetStatusShopping(ctx, conversationID); err != nil {
		return err
	}
	return nil
}

func (s *ShoppingImpl) DeleteItems(ctx context.Context, conversationID model.ConversationID, ids []string) error {
	if err := s.conversation.DeleteShoppingItems(ctx, conversationID, ids); err != nil {
		return xerrors.Errorf("failed to delete shopping item: %w", err)
	}
	return nil
}

func (s *ShoppingImpl) GetStatus(ctx context.Context, conversationID model.ConversationID) (*model.ConversationStatus, error) {
	status, err := s.conversation.GetStatus(ctx, conversationID)
	if err != nil {
		return nil, xerrors.Errorf("failed to get status: %w", err)
	}
	return status, nil
}

func (s *ShoppingImpl) SetStatus(ctx context.Context, status *model.ConversationStatus) error {
	if err := s.conversation.SetStatus(ctx, status); err != nil {
		return xerrors.Errorf("failed to set status: %w", err)
	}
	return nil
}

func (s *ShoppingImpl) SetStatusShopping(ctx context.Context, conversationID model.ConversationID) error {
	status := &model.ConversationStatus{
		ConversationID: conversationID,
		Type:           model.ConversationStatusTypeShopping,
	}
	if err := s.SetStatus(ctx, status); err != nil {
		return err
	}
	return nil
}
