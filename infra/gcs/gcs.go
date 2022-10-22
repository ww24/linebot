package gcs

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/google/wire"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
)

// Set provides a wire set.
var Set = wire.NewSet(
	New,
	NewWeatherImageStore,
	wire.Bind(new(repository.WeatherImageStore), new(*WeatherImageStore)),
	NewImageStore,
	wire.Bind(new(repository.ImageStore), new(*ImageStore)),
)

type Client struct {
	cli *storage.Client
}

func New(ctx context.Context) (*Client, error) {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to create storage client: %w", err)
	}

	return &Client{cli: cli}, nil
}
