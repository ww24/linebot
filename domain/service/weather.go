package service

import (
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/logger"
)

type Weather interface {
	SaveImage(context.Context) error
	ImageURL(context.Context) (string, error)
}

type WeatherImpl struct {
	weather    repository.Weather
	imageStore repository.WeatherImageStore
}

func NewWeather(weather repository.Weather, imageStore repository.WeatherImageStore) *WeatherImpl {
	return &WeatherImpl{
		weather:    weather,
		imageStore: imageStore,
	}
}

func (w *WeatherImpl) SaveImage(ctx context.Context) error {
	rc, err := w.weather.Fetch(ctx)
	if err != nil {
		return xerrors.Errorf("weather.Fetch: %w", err)
	}
	defer rc.Close()

	now := time.Now()
	imageURL, err := w.imageStore.Save(ctx, rc, now)
	if err != nil {
		return xerrors.Errorf("imageStore.Save: %w", err)
	}

	dl := logger.DefaultLogger(ctx)
	dl.Info("weather image saved", zap.String("url", imageURL))

	return nil
}

func (w *WeatherImpl) ImageURL(ctx context.Context) (string, error) {
	now := time.Now()
	imageURL, err := w.imageStore.Get(ctx, now)
	if err != nil {
		return "", xerrors.Errorf("imageStore.Get: %w", err)
	}
	return imageURL, nil
}
