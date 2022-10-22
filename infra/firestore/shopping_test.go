package firestore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ww24/linebot/domain/model"
)

func TestShopping_Add(t *testing.T) {
	t.Parallel()
	const conversationIDPrefix = "TestShopping_Add_"
	ctx := context.Background()
	conv := NewConversation(testCli)
	s := NewShopping(conv)
	data := []*model.ShoppingItem{
		{
			ID:             "item_01",
			Name:           "item 01",
			Quantity:       1,
			ConversationID: conversationIDPrefix + "02",
			CreatedAt:      1666416720,
			Order:          0,
		},
	}
	require.NoError(t, s.Add(ctx, data...))
	tests := []struct {
		name    string
		items   []*model.ShoppingItem
		want    []*ShoppingItem
		wantErr error
	}{
		{
			name: "create new items",
			items: []*model.ShoppingItem{
				{
					ID:             "item_01",
					Name:           "item 01",
					Quantity:       1,
					ConversationID: conversationIDPrefix + "01",
					CreatedAt:      1666416727,
					Order:          0,
				},
				{
					ID:             "item_02",
					Name:           "item 02",
					Quantity:       1,
					ConversationID: conversationIDPrefix + "01",
					CreatedAt:      1666416727,
					Order:          1,
				},
			},
			want: []*ShoppingItem{
				{
					Name:      "item 01",
					Quantity:  1,
					CreatedAt: 1666416727,
					Order:     0,
				},
				{
					Name:      "item 02",
					Quantity:  1,
					CreatedAt: 1666416727,
					Order:     1,
				},
			},
			wantErr: nil,
		},
		{
			name: "add new items",
			items: []*model.ShoppingItem{
				{
					ID:             "item_02",
					Name:           "item 02",
					Quantity:       1,
					ConversationID: conversationIDPrefix + "02",
					CreatedAt:      1666416727,
					Order:          1,
				},
				{
					ID:             "item_03",
					Name:           "item 03",
					Quantity:       1,
					ConversationID: conversationIDPrefix + "02",
					CreatedAt:      1666416727,
					Order:          2,
				},
			},
			want: []*ShoppingItem{
				{
					Name:      "item 01",
					Quantity:  1,
					CreatedAt: 1666416720,
					Order:     0,
				},
				{
					Name:      "item 02",
					Quantity:  1,
					CreatedAt: 1666416727,
					Order:     1,
				},
				{
					Name:      "item 03",
					Quantity:  1,
					CreatedAt: 1666416727,
					Order:     2,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := s.Add(ctx, tt.items...)
			require.ErrorIs(t, err, tt.wantErr)

			ss, err := s.shopping(tt.items[0].ConversationID).Documents(ctx).GetAll()
			require.NoError(t, err)
			got := make([]*ShoppingItem, 0, len(ss))
			for _, doc := range ss {
				item := new(ShoppingItem)
				require.NoError(t, doc.DataTo(item))
				got = append(got, item)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestShopping_Find(t *testing.T) {
	t.Parallel()
	const conversationID = "TestShopping_Find"
	ctx := context.Background()
	conv := NewConversation(testCli)
	s := NewShopping(conv)
	data := []*model.ShoppingItem{
		{
			ID:             "item_02",
			Name:           "item 02",
			Quantity:       2,
			ConversationID: conversationID,
			CreatedAt:      1666416727,
			Order:          1,
		},
		{
			ID:             "item_01",
			Name:           "item 01",
			Quantity:       1,
			ConversationID: conversationID,
			CreatedAt:      1666416720,
			Order:          0,
		},
	}
	require.NoError(t, s.Add(ctx, data...))
	tests := []struct {
		name           string
		conversationID model.ConversationID
		want           []*model.ShoppingItem
		wantErr        error
	}{
		{
			name:           "find all items",
			conversationID: conversationID,
			want: []*model.ShoppingItem{
				{
					ID:             "item_01",
					Name:           "item 01",
					Quantity:       1,
					ConversationID: conversationID,
					CreatedAt:      1666416720,
					Order:          0,
				},
				{
					ID:             "item_02",
					Name:           "item 02",
					Quantity:       2,
					ConversationID: conversationID,
					CreatedAt:      1666416727,
					Order:          1,
				},
			},
			wantErr: nil,
		},
		{
			name:           "no items",
			conversationID: conversationID + "_not_found",
			want:           []*model.ShoppingItem{},
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := s.Find(ctx, tt.conversationID)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestShopping_BatchDelete(t *testing.T) {
	t.Parallel()
	const conversationID = "TestShopping_BatchDelete"
	ctx := context.Background()
	conv := NewConversation(testCli)
	s := NewShopping(conv)
	data := []*model.ShoppingItem{
		{
			ID:             "item_02",
			Name:           "item 02",
			Quantity:       2,
			ConversationID: conversationID,
			CreatedAt:      1666416727,
			Order:          1,
		},
		{
			ID:             "item_01",
			Name:           "item 01",
			Quantity:       1,
			ConversationID: conversationID,
			CreatedAt:      1666416720,
			Order:          0,
		},
	}
	require.NoError(t, s.Add(ctx, data...))
	tests := []struct {
		name    string
		ids     []string
		wantErr error
	}{
		{
			name:    "delete items",
			ids:     []string{"item_01", "item_02"},
			wantErr: nil,
		},
		{
			name:    "delete not found item",
			ids:     []string{"not_found_id"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := s.BatchDelete(ctx, conversationID, tt.ids)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestShopping_DeleteAll(t *testing.T) {
	t.Parallel()
	const conversationID = "TestShopping_DeleteAll"
	ctx := context.Background()
	conv := NewConversation(testCli)
	s := NewShopping(conv)
	data := []*model.ShoppingItem{
		{
			ID:             "item_02",
			Name:           "item 02",
			Quantity:       2,
			ConversationID: conversationID,
			CreatedAt:      1666416727,
			Order:          1,
		},
		{
			ID:             "item_01",
			Name:           "item 01",
			Quantity:       1,
			ConversationID: conversationID,
			CreatedAt:      1666416720,
			Order:          0,
		},
	}
	require.NoError(t, s.Add(ctx, data...))
	tests := []struct {
		name           string
		conversationID model.ConversationID
		wantErr        error
	}{
		{
			name:           "delete all items",
			conversationID: conversationID,
			wantErr:        nil,
		},
		{
			name:           "not found",
			conversationID: conversationID + "_not_found",
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := s.DeleteAll(ctx, tt.conversationID)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
