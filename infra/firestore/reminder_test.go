package firestore

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/internal/code"
)

func TestReminder_Add(t *testing.T) {
	t.Parallel()
	const conversationID = "TestReminder_Add"
	testTime := time.Unix(1666416720, 0)
	cli := testCli.clone()
	cli.now = func() time.Time { return testTime }
	conv := NewConversation(cli)
	r := NewReminder(conv)
	ctx := context.Background()
	tests := []struct {
		name    string
		item    *model.ReminderItem
		want    *ReminderItem
		wantErr error
	}{
		{
			name: "add an item",
			item: &model.ReminderItem{
				ID:             "item_01",
				ConversationID: conversationID,
				Scheduler: &model.OneshotScheduler{
					Time: time.Unix(1666416727, 0),
				},
				Executor: &model.Executor{
					Type: model.ExecutorTypeShoppingList,
				},
			},
			want: &ReminderItem{
				Scheduler: (&model.OneshotScheduler{
					Time: time.Unix(1666416727, 0),
				}).String(),
				Executor: &Executor{
					Type: model.ExecutorTypeShoppingList,
				},
				CreatedAt: testTime.Unix(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := r.Add(ctx, tt.item)
			require.ErrorIs(t, err, tt.wantErr)

			doc, err := r.reminder(tt.item.ConversationID).Doc(string(tt.item.ID)).Get(ctx)
			require.NoError(t, err)
			got := new(ReminderItem)
			require.NoError(t, doc.DataTo(got))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReminder_List(t *testing.T) {
	t.Parallel()
	const conversationID = "TestReminder_List"
	ctx := context.Background()
	testTime := time.Unix(1666416720, 0)
	cli := testCli.clone()
	cli.now = func() time.Time { return testTime }
	conv := NewConversation(cli)
	r := NewReminder(conv)
	data := []*model.ReminderItem{
		{
			ID:             "item_01",
			ConversationID: conversationID,
			Scheduler: &model.OneshotScheduler{
				Time: time.Unix(1666416727, 0).In(time.UTC),
			},
			Executor: &model.Executor{
				Type: model.ExecutorTypeShoppingList,
			},
		},
		{
			ID:             "item_02",
			ConversationID: conversationID,
			Scheduler: &model.OneshotScheduler{
				Time: time.Unix(1666416737, 0).In(time.UTC),
			},
			Executor: &model.Executor{
				Type: model.ExecutorTypeShoppingList,
			},
		},
	}
	for _, d := range data {
		require.NoError(t, r.Add(ctx, d))
	}
	tests := []struct {
		name           string
		conversationID model.ConversationID
		want           []*model.ReminderItem
		wantErr        error
	}{
		{
			name:           "list items",
			conversationID: conversationID,
			want: []*model.ReminderItem{
				{
					ID:             "item_01",
					ConversationID: conversationID,
					Scheduler: &model.OneshotScheduler{
						Time: time.Unix(1666416727, 0).In(time.UTC),
					},
					Executor: &model.Executor{
						Type: model.ExecutorTypeShoppingList,
					},
				},
				{
					ID:             "item_02",
					ConversationID: conversationID,
					Scheduler: &model.OneshotScheduler{
						Time: time.Unix(1666416737, 0).In(time.UTC),
					},
					Executor: &model.Executor{
						Type: model.ExecutorTypeShoppingList,
					},
				},
			},
			wantErr: nil,
		},
		{
			name:           "no items",
			conversationID: conversationID + "_not_found",
			want:           []*model.ReminderItem{},
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := r.List(ctx, tt.conversationID)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReminder_Get(t *testing.T) {
	t.Parallel()
	const conversationID = "TestReminder_Get"
	testTime := time.Unix(1666416720, 0)
	cli := testCli.clone()
	cli.now = func() time.Time { return testTime }
	conv := NewConversation(cli)
	r := NewReminder(conv)
	ctx := context.Background()
	data := &model.ReminderItem{
		ID:             "item_01",
		ConversationID: conversationID,
		Scheduler: &model.OneshotScheduler{
			Time: time.Unix(1666416727, 0).In(time.UTC),
		},
		Executor: &model.Executor{
			Type: model.ExecutorTypeShoppingList,
		},
	}
	require.NoError(t, r.Add(ctx, data))
	tests := []struct {
		name     string
		itemID   model.ReminderItemID
		want     *model.ReminderItem
		wantCode code.Code
	}{
		{
			name:   "get an item",
			itemID: "item_01",
			want: &model.ReminderItem{
				ID:             "item_01",
				ConversationID: conversationID,
				Scheduler: &model.OneshotScheduler{
					Time: time.Unix(1666416727, 0).In(time.UTC),
				},
				Executor: &model.Executor{
					Type: model.ExecutorTypeShoppingList,
				},
			},
			wantCode: code.OK,
		},
		{
			name:     "not found",
			itemID:   "not_found_id",
			want:     nil,
			wantCode: code.NotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := r.Get(ctx, conversationID, tt.itemID)
			require.Equal(t, tt.wantCode, code.From(err))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReminder_Delete(t *testing.T) {
	t.Parallel()
	const conversationID = "TestReminder_Delete"
	conv := NewConversation(testCli)
	r := NewReminder(conv)
	ctx := context.Background()
	data := &model.ReminderItem{
		ID:             "item_01",
		ConversationID: conversationID,
		Scheduler: &model.OneshotScheduler{
			Time: time.Unix(1666416727, 0),
		},
		Executor: &model.Executor{
			Type: model.ExecutorTypeShoppingList,
		},
	}
	require.NoError(t, r.Add(ctx, data))
	tests := []struct {
		name    string
		itemID  model.ReminderItemID
		wantErr error
	}{
		{
			name:    "delete an item",
			itemID:  "item_01",
			wantErr: nil,
		},
		{
			name:    "not found",
			itemID:  "not_found_id",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := r.Delete(ctx, conversationID, tt.itemID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestReminder_ListAll(t *testing.T) {
	t.Parallel()
	const conversationIDPrefix = "TestReminder_ListAll_"
	conv := NewConversation(testCli)
	r := NewReminder(conv)
	ctx := context.Background()
	data := []*model.ReminderItem{
		{
			ID:             "item_01",
			ConversationID: conversationIDPrefix + "01",
			Scheduler: &model.OneshotScheduler{
				Time: time.Unix(1666416727, 0),
			},
			Executor: &model.Executor{
				Type: model.ExecutorTypeShoppingList,
			},
		},
		{
			ID:             "item_01",
			ConversationID: conversationIDPrefix + "02",
			Scheduler: &model.OneshotScheduler{
				Time: time.Unix(1666416737, 0),
			},
			Executor: &model.Executor{
				Type: model.ExecutorTypeShoppingList,
			},
		},
	}
	for _, d := range data {
		require.NoError(t, r.Add(ctx, d))
	}
	tests := []struct {
		name    string
		want    []*model.ReminderItem
		wantErr error
	}{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reminders, err := r.ListAll(ctx)
			require.ErrorIs(t, err, tt.wantErr)
			got := make([]*model.ReminderItem, 0, len(tt.want))
			for _, reminder := range reminders {
				if strings.HasPrefix(string(reminder.ConversationID), conversationIDPrefix) {
					got = append(got, reminder)
				}
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
