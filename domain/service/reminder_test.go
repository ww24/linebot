package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tenntenn/testtime"
	"go.uber.org/mock/gomock"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/mock/mock_repository"
)

func TestReminderImpl_SyncSchedule(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	items := model.ReminderItems{
		{
			ConversationID: model.ConversationID("c1"),
			Scheduler: &model.OneshotScheduler{
				Time: testTime.Add(time.Minute),
			},
		},
		{
			ConversationID: model.ConversationID("c1"),
			Scheduler: &model.OneshotScheduler{
				Time: testTime.Add(30 * time.Minute),
			},
		},
		{
			ConversationID: model.ConversationID("c2"),
			Scheduler: &model.OneshotScheduler{
				Time: testTime.Add(2 * time.Minute),
			},
		},
		{
			ConversationID: model.ConversationID("c2"),
			Scheduler: &model.OneshotScheduler{
				Time: testTime.Add(10 * time.Minute),
			},
		},
		{
			ConversationID: model.ConversationID("c3"),
			Scheduler: &model.OneshotScheduler{
				Time: testTime.Add(time.Hour),
			},
		},
	}
	tests := []struct {
		name  string
		setup func(*mock_repository.MockScheduleSynchronizer)
		items model.ReminderItems
	}{
		{
			name:  "empty",
			setup: func(m *mock_repository.MockScheduleSynchronizer) {},
			items: model.ReminderItems{},
		},
		{
			name: "one item",
			setup: func(m *mock_repository.MockScheduleSynchronizer) {
				m.EXPECT().Sync(gomock.Any(), model.ConversationID("c1"), items[0:1], testTime).Return(nil)
			},
			items: items[0:1],
		},
		{
			name: "some items",
			setup: func(m *mock_repository.MockScheduleSynchronizer) {
				m.EXPECT().Sync(gomock.Any(), model.ConversationID("c1"), items[0:2], testTime).Return(nil)
				m.EXPECT().Sync(gomock.Any(), model.ConversationID("c2"), items[2:4], testTime).Return(nil)
				m.EXPECT().Sync(gomock.Any(), model.ConversationID("c3"), items[4:5], testTime).Return(nil)
			},
			items: items[0:5],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.True(t, testtime.SetTime(t, testTime))

			ctrl := gomock.NewController(t)
			m := mock_repository.NewMockScheduleSynchronizer(ctrl)
			tt.setup(m)
			service := &ReminderImpl{
				scheduler: m,
			}

			err := service.SyncSchedule(ctx, tt.items)
			require.NoError(t, err)
		})
	}
}
