package interactor

import (
	"context"
	"io"
	"net/url"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
)

type Screenshot struct {
	browser repository.Browser
}

func NewScreenshot(browser repository.Browser) *Screenshot {
	return &Screenshot{
		browser: browser,
	}
}

func (r *Screenshot) Handle(ctx context.Context, target *url.URL, targetSelector string) (io.Reader, int, error) {
	img, size, err := r.browser.Screenshot(ctx, target, targetSelector)
	if err != nil {
		return nil, 0, xerrors.Errorf("browser.Screenshot: %w", err)
	}
	return img, size, nil
}
