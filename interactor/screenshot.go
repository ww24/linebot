package interactor

import (
	"context"
	"net/url"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/domain/service"
)

type Screenshot struct {
	browser repository.Browser
	weather service.Weather
}

func NewScreenshot(browser repository.Browser, weather service.Weather) *Screenshot {
	return &Screenshot{
		browser: browser,
		weather: weather,
	}
}

func (r *Screenshot) Handle(ctx context.Context, target *url.URL, targetSelector string) error {
	img, _, err := r.browser.Screenshot(ctx, target, targetSelector)
	if err != nil {
		return nil
	}

	if err := r.weather.SaveImage(ctx, img); err != nil {
		return xerrors.Errorf("interactor: %w", err)
	}

	return nil
}
