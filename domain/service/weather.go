package service

import (
	"context"
	"io"
	"time"

	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/internal/code"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/logger"
)

const (
	weatherImageTTL = 2 * time.Hour
	urlPathPrefix   = "/image"
)

type Weather interface {
	SaveImage(context.Context, io.Reader) error
	LatestImage(context.Context) (string, error)
}

type WeatherImpl struct {
	imageStore repository.WeatherImageStore
	loc        *time.Location
	urlPrefix  string
}

func NewWeather(
	imageStore repository.WeatherImageStore,
	ct *config.Time,
	conf *config.ServiceEndpoint,
) (*WeatherImpl, error) {
	weather := &WeatherImpl{
		imageStore: imageStore,
		loc:        ct.DefaultLocation(),
	}
	if conf.Valid() {
		endpoint, err := conf.ResolveServiceEndpoint(urlPathPrefix)
		if err != nil {
			return nil, xerrors.Errorf("failed to get endpoint: %w", err)
		}
		weather.urlPrefix = endpoint.String()
	}
	return weather, nil
}

func (w *WeatherImpl) SaveImage(ctx context.Context, r io.Reader) error {
	ctx, span := tracer.Start(ctx, "Weather#SaveImage")
	defer span.End()

	now := time.Now()
	name, err := w.imageStore.Save(ctx, r, now)
	if err != nil {
		return xerrors.Errorf("imageStore.Save: %w", err)
	}

	dl := logger.Default(ctx)
	dl.Info("service: weather image saved", zap.String("name", name))

	return nil
}

func (w *WeatherImpl) LatestImage(ctx context.Context) (string, error) {
	ctx, span := tracer.Start(ctx, "Weather#LatestImage")
	defer span.End()

	now := time.Now().In(w.loc)

	name, err := w.imageStore.Get(ctx, now, weatherImageTTL)
	if code.From(err) == code.NotFound && now.Add(-weatherImageTTL).Day() != now.Day() {
		name, err = w.imageStore.Get(ctx, now.Add(-weatherImageTTL), weatherImageTTL)
	}
	if err != nil {
		return "", xerrors.Errorf("imageStore.Get: %w", err)
	}

	imageURL := w.urlPrefix + "/" + name
	return imageURL, nil
}
