package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOneshotScheduler_Next(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scheduler *OneshotScheduler
		now       time.Time
		want      time.Time
		wantErr   error
	}{
		{
			scheduler: &OneshotScheduler{
				Time: time.Date(2021, 7, 7, 23, 59, 59, 0, time.Local),
			},
			now:  time.Date(2021, 7, 1, 0, 0, 0, 0, time.Local),
			want: time.Date(2021, 7, 7, 23, 59, 59, 0, time.Local),
		},
		{
			scheduler: &OneshotScheduler{
				Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			now:     time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC),
			want:    time.Time{},
			wantErr: ErrEndSchedule,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got, err := tt.scheduler.Next(tt.now)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDailyScheduler_Next(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scheduler *DailyScheduler
		now       time.Time
		want      time.Time
	}{
		{
			scheduler: &DailyScheduler{
				Time: time.Date(2000, 1, 1, 18, 15, 30, 0, time.Local),
			},
			now:  time.Date(2021, 4, 1, 18, 0, 0, 0, time.Local),
			want: time.Date(2021, 4, 1, 18, 15, 30, 0, time.Local),
		},
		{
			scheduler: &DailyScheduler{
				Time: time.Date(2000, 1, 1, 18, 15, 30, 0, time.UTC),
			},
			now:  time.Date(2021, 4, 1, 18, 15, 30, 0, time.UTC),
			want: time.Date(2021, 4, 2, 18, 15, 30, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got, err := tt.scheduler.Next(tt.now)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseScheduler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		serialized string
		want       Scheduler
		wantErr    error
	}{
		{
			serialized: "",
			wantErr:    ErrInvalidSchedulerType,
		},
		{
			serialized: "o#2021-04-01T10:23:45Z",
			want: &OneshotScheduler{
				Time: time.Date(2021, 4, 1, 10, 23, 45, 0, time.UTC),
			},
		},
		{
			serialized: "d#2021-01-01T15:00:00Z",
			want: &DailyScheduler{
				Time: time.Date(2021, 1, 1, 15, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got, err := ParseScheduler(tt.serialized)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
