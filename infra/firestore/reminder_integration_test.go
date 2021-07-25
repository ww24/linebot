//+build integration

package firestore

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenntenn/testtime"
	"github.com/ww24/linebot/domain/model"
)

func TestReminder_Add(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	const conversationID = "reminder-add"
	testTime := time.Now()

	// cleanup
	err := cleanupCollection(ctx, "conversations/"+conversationID+"/reminders")
	require.NoError(t, err)

	tests := []struct {
		name    string
		item    *model.ReminderItem
		want    *ReminderItem
		wantErr error
	}{
		{
			name: "add an item",
			item: &model.ReminderItem{
				Name:           "add_an_item",
				ConversationID: model.ConversationID(conversationID),
				Scheduler: &model.OnetimeScheduler{
					Time: time.Unix(10, 0),
				},
				Executor: &model.Executor{
					Type: model.ExecutorTypeShoppingList,
				},
			},
			want: &ReminderItem{
				Name:           "add_an_item",
				ConversationID: conversationID,
				Scheduler: (&model.OnetimeScheduler{
					Time: time.Unix(10, 0),
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
			require.True(t, testtime.SetTime(t, testTime))

			r := NewReminder(NewConversation(testClient))
			assert.ErrorIs(t, r.Add(ctx, tt.item), tt.wantErr)

			iter := testClient.cli.Collection("conversations/reminder-add/reminders").
				Where("name", "==", tt.item.Name).Documents(ctx)
			docs, err := iter.GetAll()
			require.NoError(t, err)
			require.Len(t, docs, 1)

			got := new(ReminderItem)
			require.NoError(t, docs[0].DataTo(got))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReminder_Find(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	const conversationID = "reminder-find-a"

	// cleanup
	err := cleanupCollection(ctx, "conversations/"+conversationID+"/reminders")
	require.NoError(t, err)

	// setup
	batch := testClient.cli.Batch()
	colRef := testClient.cli.Collection("conversations/" + conversationID + "/reminders")
	batch.Create(colRef.Doc("reminder-find-a#test01"), &ReminderItem{
		Name:           "test01",
		ConversationID: conversationID,
		Scheduler: (&model.OnetimeScheduler{
			Time: time.Unix(11, 0),
		}).String(),
		Executor: &Executor{
			Type: model.ExecutorTypeShoppingList,
		},
		CreatedAt: 1,
	})
	batch.Create(colRef.Doc("reminder-find-a#test02"), &ReminderItem{
		Name:           "test02",
		ConversationID: conversationID,
		Scheduler: (&model.DailyScheduler{
			Time: time.Unix(12, 0),
		}).String(),
		Executor: &Executor{
			Type: model.ExecutorTypeShoppingList,
		},
		CreatedAt: 2,
	})
	_, err = batch.Commit(ctx)
	require.NoError(t, err)

	tests := []struct {
		name           string
		conversationID model.ConversationID
		want           []*model.ReminderItem
		wantErr        error
	}{
		{
			name:           "find reminders",
			conversationID: conversationID,
			want: []*model.ReminderItem{
				{
					ID:             "reminder-find-a#test01",
					Name:           "test01",
					ConversationID: conversationID,
					Scheduler: &model.OnetimeScheduler{
						Time: time.Unix(11, 0),
					},
					Executor: &model.Executor{
						Type: model.ExecutorTypeShoppingList,
					},
				},
				{
					ID:             "reminder-find-a#test02",
					Name:           "test02",
					ConversationID: conversationID,
					Scheduler: &model.DailyScheduler{
						Time: time.Unix(12, 0),
					},
					Executor: &model.Executor{
						Type: model.ExecutorTypeShoppingList,
					},
				},
			},
			wantErr: nil,
		},
		{
			name:           "find reminders with empty collection",
			conversationID: "reminder-find-b",
			want:           []*model.ReminderItem{},
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := NewReminder(NewConversation(testClient))
			got, err := r.Find(ctx, tt.conversationID)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReminder_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	const conversationID = "reminder-delete-a"

	// cleanup
	err := cleanupCollection(ctx, "conversations/"+conversationID+"/reminders")
	require.NoError(t, err)

	// setup
	batch := testClient.cli.Batch()
	colRef := testClient.cli.Collection("conversations/" + conversationID + "/reminders")
	batch.Create(colRef.Doc("reminder-delete-a#test01"), &ReminderItem{
		Name:           "test01",
		ConversationID: conversationID,
		Scheduler: (&model.OnetimeScheduler{
			Time: time.Unix(11, 0),
		}).String(),
		Executor: &Executor{
			Type: model.ExecutorTypeShoppingList,
		},
		CreatedAt: 1,
	})
	batch.Create(colRef.Doc("reminder-delete-a#test02"), &ReminderItem{
		Name:           "test02",
		ConversationID: conversationID,
		Scheduler: (&model.DailyScheduler{
			Time: time.Unix(12, 0),
		}).String(),
		Executor: &Executor{
			Type: model.ExecutorTypeShoppingList,
		},
		CreatedAt: 2,
	})
	_, err = batch.Commit(ctx)
	require.NoError(t, err)

	tests := []struct {
		name           string
		conversationID model.ConversationID
		id             string
		want           []*ReminderItem
		wantErr        error
	}{
		{
			name:           "delete a reminder",
			conversationID: conversationID,
			id:             "reminder-delete-a#test01",
			want: []*ReminderItem{
				{
					Name:           "test02",
					ConversationID: conversationID,
					Scheduler: (&model.DailyScheduler{
						Time: time.Unix(12, 0),
					}).String(),
					Executor: &Executor{
						Type: model.ExecutorTypeShoppingList,
					},
					CreatedAt: 2,
				},
			},
			wantErr: nil,
		},
		{
			name:           "delete not exist item",
			conversationID: "reminder-delete-b",
			id:             "reminder-delete-b#test01",
			want:           []*ReminderItem{},
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewReminder(NewConversation(testClient))
			assert.ErrorIs(t, r.Delete(ctx, tt.conversationID, tt.id), tt.wantErr)

			iter := testClient.cli.Collection("conversations/" + string(tt.conversationID) + "/reminders").
				Documents(ctx)
			docs, err := iter.GetAll()
			require.NoError(t, err)

			got := make([]*ReminderItem, 0, len(docs))
			for _, doc := range docs {
				item := &ReminderItem{}
				require.NoError(t, doc.DataTo(item))
				got = append(got, item)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
