package interactor

import (
	"context"
	"io"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
)

type Image struct {
	imageStore repository.ImageStore
}

func NewImage(imageStore repository.ImageStore) *Image {
	return &Image{
		imageStore: imageStore,
	}
}

func (i *Image) Handle(ctx context.Context, key string) (io.ReadCloser, int, error) {
	rc, size, err := i.imageStore.Fetch(ctx, key)
	if err != nil {
		return nil, 0, xerrors.Errorf("imageStore.Fetch: %w", err)
	}
	return rc, size, nil
}
