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

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/logger"
)

const (
	weatherPrefix   = "weather/japan-all/"
	objectSuffix    = "-weather.png"
	urlPathPrefix   = "/image"
	weatherImageTTL = 2 * time.Hour
)

type WeatherImageStore struct {
	*Client
	bucket    string
	loc       *time.Location
	urlPrefix string
}

func NewWeatherImageStore(cli *Client, conf repository.Config) (*WeatherImageStore, error) {
	endpoint, err := conf.ServiceEndpoint(urlPathPrefix)
	if err != nil {
		return nil, xerrors.Errorf("failed to get endpoint: %w", err)
	}

	return &WeatherImageStore{
		Client:    cli,
		bucket:    conf.ImageBucket(),
		loc:       conf.DefaultLocation(),
		urlPrefix: endpoint.String(),
	}, nil
}

func (w *WeatherImageStore) Save(ctx context.Context, r io.Reader, t time.Time) (string, error) {
	key := w.key(t)
	obj := w.cli.Bucket(w.bucket).Object(key)
	writer := obj.NewWriter(ctx)

	dl := logger.DefaultLogger(ctx)
	dl.Info("upload image", zap.String("key", key))

	if _, err := io.Copy(writer, r); err != nil {
		return "", xerrors.Errorf("io.Copy: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", xerrors.Errorf("writer.Close: %w", err)
	}

	return w.url(key), nil
}

func (w *WeatherImageStore) Get(ctx context.Context, t time.Time) (string, error) {
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
		if attrs.Created.Add(weatherImageTTL).Before(t) {
			return "", xerrors.Errorf("image is expired")
		}

		return w.url(attrs.Name), nil
	}

	return "", xerrors.Errorf("image is not found")
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

func (w *WeatherImageStore) url(key string) string {
	return w.urlPrefix + "/" + key
}
