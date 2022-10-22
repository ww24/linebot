package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ww24/linebot/domain/model"
)

type Shopping struct {
	*Conversation
}

func NewShopping(c *Conversation) *Shopping {
	return &Shopping{Conversation: c}
}

func (s *Shopping) shopping(conversationID model.ConversationID) *firestore.CollectionRef {
	return s.conversation(conversationID).Collection("shoppings")
}

func (s *Shopping) Add(ctx context.Context, items ...*model.ShoppingItem) error {
	ctx, span := tracer.Start(ctx, "Shopping#Add")
	defer span.End()

	txf := func(ctx context.Context, tx *firestore.Transaction) error {
		for _, item := range items {
			if err := item.Validate(); err != nil {
				return xerrors.Errorf("shopping item validation failed: %w", err)
			}

			entity := NewShoppingItem(item)
			shopping := s.shopping(item.ConversationID)
			if err := tx.Create(shopping.Doc(entity.ID), entity); err != nil {
				return xerrors.Errorf("failed to create: %w", err)
			}
		}
		return nil
	}
	if err := s.cli.RunTransaction(ctx, txf); err != nil {
		return xerrors.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (s *Shopping) Find(ctx context.Context, conversationID model.ConversationID) ([]*model.ShoppingItem, error) {
	ctx, span := tracer.Start(ctx, "Shopping#Find")
	defer span.End()

	iter := s.shopping(conversationID).
		OrderBy("created_at", firestore.Asc).
		OrderBy("order", firestore.Asc).
		Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, xerrors.Errorf("failed to get all: %w", err)
	}

	items := make([]*model.ShoppingItem, 0, len(docs))
	for _, doc := range docs {
		var item ShoppingItem
		if err := doc.DataTo(&item); err != nil {
			return nil, xerrors.Errorf("failed to convert response as ShoppingItem: %w", err)
		}
		items = append(items, item.Model(conversationID, doc.Ref.ID))
	}

	return items, nil
}

func (s *Shopping) BatchDelete(ctx context.Context, conversationID model.ConversationID, ids []string) error {
	ctx, span := tracer.Start(ctx, "Shopping#BatchDelete")
	defer span.End()

	txf := func(ctx context.Context, tx *firestore.Transaction) error {
		for _, id := range ids {
			item := s.shopping(conversationID).Doc(id)
			if err := tx.Delete(item, firestore.Exists); err != nil {
				return xerrors.Errorf("failed to delete document: %w", err)
			}
		}
		return nil
	}
	if err := s.cli.RunTransaction(ctx, txf); err != nil {
		if status.Code(err) != codes.NotFound {
			return xerrors.Errorf("transaction failed: %w", err)
		}
	}

	return nil
}

func (s *Shopping) DeleteAll(ctx context.Context, conversationID model.ConversationID) error {
	ctx, span := tracer.Start(ctx, "Shopping#DeleteAll")
	defer span.End()

	refs, err := s.shopping(conversationID).DocumentRefs(ctx).GetAll()
	if err != nil {
		return xerrors.Errorf("failed to get document refs: %w", err)
	}
	if len(refs) == 0 {
		return nil
	}

	txf := func(ctx context.Context, tx *firestore.Transaction) error {
		for _, ref := range refs {
			if err := tx.Delete(ref, firestore.Exists); err != nil {
				return xerrors.Errorf("failed to delete document: %w", err)
			}
		}
		return nil
	}
	if err := s.cli.RunTransaction(ctx, txf); err != nil {
		return xerrors.Errorf("transaction failed: %w", err)
	}

	return nil
}

type ShoppingItem struct {
	ConversationID model.ConversationID `firestore:"-"`
	ID             string               `firestore:"-"`
	Name           string               `firestore:"name"`
	Quantity       int                  `firestore:"quantity"`
	CreatedAt      int64                `firestore:"created_at"`
	Order          int                  `firestore:"order"`
}

func NewShoppingItem(src *model.ShoppingItem) *ShoppingItem {
	return &ShoppingItem{
		ConversationID: src.ConversationID,
		ID:             src.ID,
		Name:           src.Name,
		Quantity:       src.Quantity,
		CreatedAt:      src.CreatedAt,
		Order:          src.Order,
	}
}

func (c *ShoppingItem) Model(conversationID model.ConversationID, id string) *model.ShoppingItem {
	return &model.ShoppingItem{
		ConversationID: conversationID,
		ID:             id,
		Name:           c.Name,
		Quantity:       c.Quantity,
		CreatedAt:      c.CreatedAt,
		Order:          c.Order,
	}
}
