package linebot

import (
	"embed"
	"testing"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ww24/linebot/domain/model"
)

var (
	//go:embed testdata/*.golden
	golden embed.FS
)

func TestMakeReminderListMessage(t *testing.T) {
	t.Parallel()
	testTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	items := []*model.ReminderItem{
		{
			ID:             "id1",
			ConversationID: "conversationID1",
			Scheduler: &model.OneshotScheduler{
				Time: time.Date(2020, 1, 3, 12, 30, 0, 0, time.UTC),
			},
			Executor: &model.Executor{
				Type: model.ExecutorTypeShoppingList,
			},
		},
	}
	want, err := golden.ReadFile("testdata/reminder_list.json.golden")
	require.NoError(t, err)
	got, err := makeReminderListMessage(items, testTime)
	require.NoError(t, err)
	assert.Equal(t, string(want), string(got))

	_, err = linebot.UnmarshalFlexMessageJSON(got)
	require.NoError(t, err, "invalid flex message")
}

func TestToReminderItem(t *testing.T) {
	t.Parallel()
	testTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		item *model.ReminderItem
		want *ReminderItem
	}{
		{
			item: &model.ReminderItem{
				ID:             "id1",
				ConversationID: "conversationID1",
				Scheduler: &model.OneshotScheduler{
					Time: time.Date(2020, 1, 3, 12, 30, 0, 0, time.UTC),
				},
				Executor: &model.Executor{
					Type: model.ExecutorTypeShoppingList,
				},
			},
			want: &ReminderItem{
				Title:        "買い物リスト",
				SubTitle:     "at 2020-01-03 12:30.",
				Next:         "01/03 12:30",
				DeleteTarget: "Reminder#delete#id1",
			},
		},
		{
			item: &model.ReminderItem{
				ID:             "id2",
				ConversationID: "conversationID1",
				Scheduler: &model.DailyScheduler{
					Time: time.Date(2020, 1, 3, 12, 30, 0, 0, time.UTC),
				},
				Executor: &model.Executor{
					Type: model.ExecutorTypeShoppingList,
				},
			},
			want: &ReminderItem{
				Title:        "買い物リスト",
				SubTitle:     "at 12:30 every day.",
				Next:         "01/01 12:30",
				DeleteTarget: "Reminder#delete#id2",
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got := toReminderItem(tt.item, testTime)
			assert.Equal(t, tt.want, got)
		})
	}
}
