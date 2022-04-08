package gcs

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/google/wire"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/xerrors"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

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
	opts := []option.ClientOption{
		option.WithGRPCDialOption(
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		),
		option.WithGRPCDialOption(
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		),
	}

	cli, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, xerrors.Errorf("failed to create storage client: %w", err)
	}

	return &Client{cli: cli}, nil
}
