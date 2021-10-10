package firestore

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"

	"github.com/ww24/linebot/domain/model"
)

type Conversation struct {
	*Client
}

type ShoppingItem struct {
	ID             string `firestore:"-"`
	Name           string `firestore:"name"`
	Quantity       int    `firestore:"quantity"`
	ConversationID string `firestore:"conversation_id"`
	CreatedAt      int64  `firestore:"created_at"`
	Order          int    `firestore:"order"`
}

func NewShoppingItem(src *model.ShoppingItem) *ShoppingItem {
	return &ShoppingItem{
		ID:             src.ID,
		Name:           src.Name,
		Quantity:       src.Quantity,
		ConversationID: string(src.ConversationID),
		CreatedAt:      src.CreatedAt,
		Order:          src.Order,
	}
}

func (c *ShoppingItem) Model() *model.ShoppingItem {
	return &model.ShoppingItem{
		ID:             c.ID,
		Name:           c.Name,
		Quantity:       c.Quantity,
		ConversationID: model.ConversationID(c.ConversationID),
		CreatedAt:      c.CreatedAt,
		Order:          c.Order,
	}
}

type ConversationStatus struct {
	ConversationID model.ConversationID
	Type           int `firestore:"type"`
}

func (c *ConversationStatus) Model() *model.ConversationStatus {
	return &model.ConversationStatus{
		ConversationID: c.ConversationID,
		Type:           model.ConversationStatusType(c.Type),
	}
}

func NewConversationStatus(src *model.ConversationStatus) *ConversationStatus {
	return &ConversationStatus{
		ConversationID: src.ConversationID,
		Type:           int(src.Type),
	}
}

func NewConversation(cli *Client) *Conversation {
	return &Conversation{Client: cli}
}

func (c *Conversation) conversation(conversationID model.ConversationID) *firestore.DocumentRef {
	return c.cli.Collection("conversations").Doc(string(conversationID))
}

func (c *Conversation) shopping(conversationID model.ConversationID) *firestore.CollectionRef {
	return c.conversation(conversationID).Collection("shoppings")
}

func (c *Conversation) AddShoppingItem(ctx context.Context, items ...*model.ShoppingItem) error {
	tracer := otel.Tracer("infra/firestore")
	ctx, span := tracer.Start(ctx, "AddShoppingItem")
	defer span.End()

	batch := c.cli.Batch()
	for _, item := range items {
		if err := item.Validate(); err != nil {
			return xerrors.Errorf("shopping item validation failed: %w", err)
		}

		entity := NewShoppingItem(item)
		shopping := c.shopping(item.ConversationID)
		batch.Create(shopping.NewDoc(), entity)
	}

	if _, err := batch.Commit(ctx); err != nil {
		return xerrors.Errorf("failed to commit: %w", err)
	}
	return nil
}

func (c *Conversation) FindShoppingItem(ctx context.Context, conversationID model.ConversationID) ([]*model.ShoppingItem, error) {
	tracer := otel.Tracer("infra/firestore")
	ctx, span := tracer.Start(ctx, "FindShoppingItem")
	defer span.End()

	iter := c.shopping(conversationID).
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
		items = append(items, item.Model())
	}

	return items, nil
}

func (c *Conversation) DeleteShoppingItems(ctx context.Context, conversationID model.ConversationID, ids []string) error {
	tracer := otel.Tracer("infra/firestore")
	ctx, span := tracer.Start(ctx, "DeleteShoppingItems")
	defer span.End()

	batch := c.cli.Batch()
	for _, id := range ids {
		item := c.shopping(conversationID).Doc(id)
		batch.Delete(item, firestore.Exists)
	}

	if _, err := batch.Commit(ctx); err != nil {
		return xerrors.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (c *Conversation) DeleteAllShoppingItem(ctx context.Context, conversationID model.ConversationID) error {
	tracer := otel.Tracer("infra/firestore")
	ctx, span := tracer.Start(ctx, "DeleteAllShoppingItem")
	defer span.End()

	iter := c.shopping(conversationID).DocumentRefs(ctx)
	batch := c.cli.Batch()

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

func (c *Conversation) SetStatus(ctx context.Context, status *model.ConversationStatus) error {
	tracer := otel.Tracer("infra/firestore")
	ctx, span := tracer.Start(ctx, "SetStatus")
	defer span.End()

	if err := status.Validate(); err != nil {
		return xerrors.Errorf("conversation status validation failed: %w", err)
	}

	dr := c.conversation(status.ConversationID).Collection("status").Doc("#")
	entity := NewConversationStatus(status)
	if _, err := dr.Set(ctx, entity); err != nil {
		return xerrors.Errorf("failed to set conversation status: %w", err)
	}

	return nil
}

func (c *Conversation) GetStatus(ctx context.Context, conversationID model.ConversationID) (*model.ConversationStatus, error) {
	tracer := otel.Tracer("infra/firestore")
	ctx, span := tracer.Start(ctx, "GetStatus")
	defer span.End()

	doc, err := c.conversation(conversationID).Collection("status").Doc("#").Get(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to get conversation status: %w", err)
	}

	var ret ConversationStatus
	if err := doc.DataTo(&ret); err != nil {
		return nil, xerrors.Errorf("failed to convert response as ConversationStatus: %w", err)
	}
	return ret.Model(), nil
}
