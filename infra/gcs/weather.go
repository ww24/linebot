package gcs

import (
	"context"
	"errors"
	"io"
	"math"
	"path"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"

	"github.com/ww24/linebot/internal/code"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/logger"
)

const (
	weatherPrefix = "weather/japan-all/"
	objectSuffix  = "-weather.png"
)

type WeatherImageStore struct {
	*Client
	bucket string
	loc    *time.Location
}

func NewWeatherImageStore(cli *Client, cs *config.Storage, ct *config.Time) (*WeatherImageStore, error) {
	return &WeatherImageStore{
		Client: cli,
		bucket: cs.ImageBucket,
		loc:    ct.DefaultLocation(),
	}, nil
}

func (w *WeatherImageStore) Save(ctx context.Context, r io.Reader, t time.Time) (string, error) {
	key := w.key(t)
	obj := w.cli.Bucket(w.bucket).Object(key)
	writer := obj.NewWriter(ctx)

	dl := logger.Default(ctx)
	dl.Info("upload image", zap.String("key", key))

	if _, err := io.Copy(writer, r); err != nil {
		return "", xerrors.Errorf("io.Copy: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", xerrors.Errorf("writer.Close: %w", err)
	}

	return key, nil
}

func (w *WeatherImageStore) Get(ctx context.Context, t time.Time, ttl time.Duration) (string, error) {
	q := &storage.Query{
		Delimiter: "/",
		Prefix:    weatherPrefix + t.In(w.loc).Format("20060102") + "/",
	}
	iter := w.cli.Bucket(w.bucket).Objects(ctx, q)
	for {
		attrs, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return "", xerrors.Errorf("failed to get image: %w", err)
		}
		if !strings.HasSuffix(attrs.Name, objectSuffix) {
			continue
		}
		if attrs.Created.Add(ttl).Before(t) {
			return "", xerrors.Errorf("image is expired")
		}

		return attrs.Name, nil
	}

	err := xerrors.Errorf("image is not found")
	return "", code.With(err, code.NotFound)
}

func (w *WeatherImageStore) key(t time.Time) string {
	reverseUnixtime := math.MaxInt64 - t.Unix()
	const base = 10
	return path.Join(
		weatherPrefix,
		t.In(w.loc).Format("20060102"),
		strconv.FormatInt(reverseUnixtime, base)+objectSuffix,
	)
}
