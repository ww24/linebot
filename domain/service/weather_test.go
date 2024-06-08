package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tenntenn/testtime"
	"go.uber.org/mock/gomock"

	"github.com/ww24/linebot/internal/code"
	"github.com/ww24/linebot/mock/mock_repository"
)

func TestWeatherImpl_LatestImage(t *testing.T) {
	t.Parallel()
	const urlPrefix = "https://example.com/image"
	ctx := context.Background()
	loc := time.FixedZone("Asia/Tokyo", 9*60*60)
	tests := []struct {
		name     string
		setup    func(*mock_repository.MockWeatherImageStore, time.Time)
		time     time.Time
		want     string
		wantCode code.Code
	}{
		{
			name: "success",
			setup: func(m *mock_repository.MockWeatherImageStore, t time.Time) {
				m.EXPECT().Get(
					gomock.Any(), t,
					weatherImageTTL,
				).Return("20220101/image.png", nil)
			},
			time:     time.Date(2022, 1, 1, 12, 0, 0, 0, loc),
			want:     urlPrefix + "/20220101/image.png",
			wantCode: code.OK,
		},
		{
			name: "today's image is not found, but yesterday's one is available",
			setup: func(m *mock_repository.MockWeatherImageStore, t time.Time) {
				gomock.InOrder(
					m.EXPECT().Get(
						gomock.Any(), t,
						weatherImageTTL,
					).Return("", code.With(errors.New("not found"), code.NotFound)),
					m.EXPECT().Get(
						gomock.Any(), t.Add(-weatherImageTTL),
						weatherImageTTL,
					).Return("20211231/image.png", nil),
				)
			},
			time:     time.Date(2022, 1, 1, 1, 0, 0, 0, loc),
			want:     urlPrefix + "/20211231/image.png",
			wantCode: code.OK,
		},
		{
			name: "not found",
			setup: func(m *mock_repository.MockWeatherImageStore, t time.Time) {
				gomock.InOrder(
					m.EXPECT().Get(
						gomock.Any(), t,
						weatherImageTTL,
					).Return("", code.With(errors.New("not found"), code.NotFound)),
					m.EXPECT().Get(
						gomock.Any(), t.Add(-weatherImageTTL),
						weatherImageTTL,
					).Return("", code.With(errors.New("not found"), code.NotFound)),
				)
			},
			time:     time.Date(2022, 1, 1, 1, 0, 0, 0, loc),
			want:     "",
			wantCode: code.NotFound,
		},
		{
			name: "unexpected error",
			setup: func(m *mock_repository.MockWeatherImageStore, t time.Time) {
				m.EXPECT().Get(
					gomock.Any(), t,
					weatherImageTTL,
				).Return("", errors.New("unexpected"))
			},
			time:     time.Date(2022, 1, 1, 1, 0, 0, 0, loc),
			want:     "",
			wantCode: code.Unexpected,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.True(t, testtime.SetTime(t, tt.time))

			ctrl := gomock.NewController(t)
			m := mock_repository.NewMockWeatherImageStore(ctrl)
			tt.setup(m, tt.time)
			service := &WeatherImpl{
				imageStore: m,
				loc:        loc,
				urlPrefix:  urlPrefix,
			}

			got, err := service.LatestImage(ctx)
			assert.Equal(t, tt.wantCode, code.From(err))
			assert.Equal(t, tt.want, got)
		})
	}
}
