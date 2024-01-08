package gcs

import (
	"context"
	"errors"
	"io"

	"cloud.google.com/go/storage"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/internal/code"
	"github.com/ww24/linebot/internal/config"
)

type ImageStore struct {
	*Client
	bucket string
}

func NewImageStore(cli *Client, conf *config.Storage) (*ImageStore, error) {
	return &ImageStore{
		Client: cli,
		bucket: conf.ImageBucket,
	}, nil
}

func (w *ImageStore) Fetch(ctx context.Context, key string) (io.ReadCloser, int, error) {
	obj := w.cli.Bucket(w.bucket).Object(key)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, 0, code.With(err, code.NotFound)
		}
		return nil, 0, xerrors.Errorf("failed to get reader: %w", err)
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to get attrs: %w", err)
	}

	return reader, int(attrs.Size), nil
}
