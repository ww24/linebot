package service

import (
	"context"
	"io"
	"time"

	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/internal/code"
	"github.com/ww24/linebot/logger"
)

const (
	weatherImageTTL = 2 * time.Hour
)

type Weather interface {
	SaveImage(context.Context, io.Reader) error
	ImageURL(context.Context) (string, error)
}

type WeatherImpl struct {
	imageStore repository.WeatherImageStore
	loc        *time.Location
}

func NewWeather(
	imageStore repository.WeatherImageStore,
	conf repository.Config,
) *WeatherImpl {
	return &WeatherImpl{
		imageStore: imageStore,
		loc:        conf.DefaultLocation(),
	}
}

func (w *WeatherImpl) SaveImage(ctx context.Context, r io.Reader) error {
	ctx, span := tracer.Start(ctx, "Weather#SaveImage")
	defer span.End()

	now := time.Now()
	imageURL, err := w.imageStore.Save(ctx, r, now)
	if err != nil {
		return xerrors.Errorf("imageStore.Save: %w", err)
	}

	dl := logger.DefaultLogger(ctx)
	dl.Info("weather image saved", zap.String("url", imageURL))

	return nil
}

func (w *WeatherImpl) ImageURL(ctx context.Context) (string, error) {
	ctx, span := tracer.Start(ctx, "Weather#ImageURL")
	defer span.End()

	now := time.Now().In(w.loc)

	imageURL, err := w.imageStore.Get(ctx, now, weatherImageTTL)
	if code.From(err) == code.NotFound && now.Add(-weatherImageTTL).Day() != now.Day() {
		imageURL, err = w.imageStore.Get(ctx, now.Add(-weatherImageTTL), weatherImageTTL)
	}
	if err != nil {
		return "", xerrors.Errorf("imageStore.Get: %w", err)
	}

	return imageURL, nil
}
