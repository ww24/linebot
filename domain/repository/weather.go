//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOPACKAGE/mock_$GOFILE -package=mock_repository

package repository

import (
	"context"
	"io"
	"time"
)

type Weather interface {
	Fetch(context.Context) (io.ReadCloser, error)
}

type WeatherImageStore interface {
	Save(context.Context, io.Reader, time.Time) (string, error)
	Get(context.Context, time.Time) (string, error)
}

type ImageStore interface {
	Fetch(context.Context, string) (io.ReadCloser, int, error)
}
