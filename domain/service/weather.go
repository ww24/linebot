package service

import (
	"context"
	"time"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
)

type Weather interface {
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

func (w *WeatherImpl) ImageURL(ctx context.Context) (string, error) {
	now := time.Now()
	imageURL, err := w.imageStore.Get(ctx, now)
	if err == nil {
		return imageURL, nil
	}

	// fetch image if not exists in store
	rc, err := w.weather.Fetch(ctx)
	if err != nil {
		return "", xerrors.Errorf("weather.Fetch: %w", err)
	}
	defer rc.Close()

	imageURL, err = w.imageStore.Save(ctx, rc, now)
	if err != nil {
		return "", xerrors.Errorf("imageStore.Save: %w", err)
	}

	return imageURL, nil
}
