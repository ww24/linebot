package repository

import (
	"context"
	"io"
	"net/url"
)

type Browser interface {
	Screenshot(context.Context, *url.URL, string) (io.Reader, int, error)
}
