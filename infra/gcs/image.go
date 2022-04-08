package gcs

import (
	"context"
	"io"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
)

type ImageStore struct {
	*Client
	bucket string
}

func NewImageStore(cli *Client, conf repository.Config) (*ImageStore, error) {
	return &ImageStore{
		Client: cli,
		bucket: conf.ImageBucket(),
	}, nil
}

func (w *ImageStore) Fetch(ctx context.Context, key string) (io.ReadCloser, int, error) {
	obj := w.cli.Bucket(w.bucket).Object(key)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to get reader: %w", err)
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, 0, xerrors.Errorf("failed to get attrs: %w", err)
	}

	return reader, int(attrs.Size), nil
}
