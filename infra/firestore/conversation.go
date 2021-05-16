package firestore

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/ww24/linebot/domain/model"
	"google.golang.org/api/iterator"
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
}

func NewShoppingItem(src *model.ShoppingItem) *ShoppingItem {
	return &ShoppingItem{
		ID:             src.ID,
		Name:           src.Name,
		Quantity:       src.Quantity,
		ConversationID: string(src.ConversationID),
		CreatedAt:      src.CreatedAt,
	}
}

func (c *ShoppingItem) Model() *model.ShoppingItem {
	return &model.ShoppingItem{
		ID:             c.ID,
		Name:           c.Name,
		Quantity:       c.Quantity,
		ConversationID: model.ConversationID(c.ConversationID),
		CreatedAt:      c.CreatedAt,
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
	batch := c.cli.Batch()
	for _, item := range items {
		if err := item.Validate(); err != nil {
			return err
		}

		entity := NewShoppingItem(item)
		shopping := c.shopping(item.ConversationID)
		batch.Create(shopping.NewDoc(), entity)
	}

	if _, err := batch.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Conversation) FindShoppingItem(ctx context.Context, conversationID model.ConversationID) ([]*model.ShoppingItem, error) {
	iter := c.shopping(conversationID).
		OrderBy("created_at", firestore.Asc).
		Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}

	items := make([]*model.ShoppingItem, 0, len(docs))
	for _, doc := range docs {
		var item ShoppingItem
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}
		item.ID = doc.Ref.ID
		items = append(items, item.Model())
	}

	return items, nil
}

func (c *Conversation) DeleteShoppingItems(ctx context.Context, conversationID model.ConversationID, ids []string) error {
	batch := c.cli.Batch()
	for _, id := range ids {
		item := c.shopping(conversationID).Doc(id)
		batch.Delete(item, firestore.Exists)
	}

	if _, err := batch.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Conversation) DeleteAllShoppingItem(ctx context.Context, conversationID model.ConversationID) error {
	iter := c.shopping(conversationID).DocumentRefs(ctx)
	batch := c.cli.Batch()

	nothing := true
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return err
		}

		batch.Delete(doc, firestore.Exists)
		nothing = false
	}

	if nothing {
		return nil
	}

	if _, err := batch.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Conversation) SetStatus(ctx context.Context, status *model.ConversationStatus) error {
	if err := status.Validate(); err != nil {
		return err
	}

	dr := c.conversation(status.ConversationID).Collection("status").Doc("#")
	entity := NewConversationStatus(status)
	if _, err := dr.Set(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (c *Conversation) GetStatus(ctx context.Context, conversationID model.ConversationID) (*model.ConversationStatus, error) {
	doc, err := c.conversation(conversationID).Collection("status").Doc("#").Get(ctx)
	if err != nil {
		return nil, err
	}

	var ret ConversationStatus
	if err := doc.DataTo(&ret); err != nil {
		return nil, err
	}
	return ret.Model(), nil
}
