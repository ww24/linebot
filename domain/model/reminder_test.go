package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReminderItems_FilterNextSchedule(t *testing.T) {
	t.Parallel()
	testTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name  string
		items ReminderItems
		d     time.Duration
		want  ReminderItems
	}{
		{
			name:  "empty",
			items: ReminderItems{},
			want:  ReminderItems{},
		},
		{
			name: "no matches",
			items: ReminderItems{
				{
					Scheduler: &OneshotScheduler{
						Time: testTime.Add(-time.Hour),
					},
				},
				{
					Scheduler: &OneshotScheduler{
						Time: testTime.Add(time.Hour),
					},
				},
				{
					Scheduler: &DailyScheduler{
						Time: testTime.Add(time.Hour),
					},
				},
				{
					Scheduler: &DailyScheduler{
						Time: testTime,
					},
				},
			},
			d:    time.Hour,
			want: ReminderItems{},
		},
		{
			name: "matches",
			items: ReminderItems{
				{
					Scheduler: &OneshotScheduler{
						Time: testTime.Add(time.Minute),
					},
				},
				{
					Scheduler: &DailyScheduler{
						Time: testTime.Add(-23*time.Hour - 30*time.Minute),
					},
				},
			},
			d: time.Hour,
			want: ReminderItems{
				{
					Scheduler: &OneshotScheduler{
						Time: testTime.Add(time.Minute),
					},
				},
				{
					Scheduler: &DailyScheduler{
						Time: testTime.Add(-23*time.Hour - 30*time.Minute),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.items.FilterNextSchedule(testTime, tt.d)
			assert.Equal(t, tt.want, got)
		})
	}
}
