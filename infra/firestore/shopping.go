package firestore

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"

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

	batch := s.cli.Batch()
	for _, item := range items {
		if err := item.Validate(); err != nil {
			return xerrors.Errorf("shopping item validation failed: %w", err)
		}

		entity := NewShoppingItem(item)
		shopping := s.shopping(item.ConversationID)
		batch.Create(shopping.NewDoc(), entity)
	}

	if _, err := batch.Commit(ctx); err != nil {
		return xerrors.Errorf("failed to commit: %w", err)
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
		item.ID = doc.Ref.ID
		item.ConversationID = conversationID
		items = append(items, item.Model())
	}

	return items, nil
}

func (s *Shopping) BatchDelete(ctx context.Context, conversationID model.ConversationID, ids []string) error {
	ctx, span := tracer.Start(ctx, "Shopping#BatchDelete")
	defer span.End()

	batch := s.cli.Batch()
	for _, id := range ids {
		item := s.shopping(conversationID).Doc(id)
		batch.Delete(item, firestore.Exists)
	}

	if _, err := batch.Commit(ctx); err != nil {
		return xerrors.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (s *Shopping) DeleteAll(ctx context.Context, conversationID model.ConversationID) error {
	ctx, span := tracer.Start(ctx, "Shopping#DeleteAll")
	defer span.End()

	iter := s.shopping(conversationID).DocumentRefs(ctx)
	batch := s.cli.Batch()

	nothing := true
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return xerrors.Errorf("failed to iterate: %w", err)
		}

		batch.Delete(doc, firestore.Exists)
		nothing = false
	}

	if nothing {
		return nil
	}

	if _, err := batch.Commit(ctx); err != nil {
		return xerrors.Errorf("failed to commit: %w", err)
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

func (c *ShoppingItem) Model() *model.ShoppingItem {
	return &model.ShoppingItem{
		ConversationID: c.ConversationID,
		ID:             c.ID,
		Name:           c.Name,
		Quantity:       c.Quantity,
		CreatedAt:      c.CreatedAt,
		Order:          c.Order,
	}
}
