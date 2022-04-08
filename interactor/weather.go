package interactor

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/domain/service"
)

const (
	triggerWeather = "天気"
)

type Weather struct {
	weather service.Weather
	message repository.MessageProviderSet
	bot     service.Bot
}

func NewWeather(
	weather service.Weather,
	message repository.MessageProviderSet,
	bot service.Bot,
) *Weather {
	return &Weather{
		weather: weather,
		message: message,
		bot:     bot,
	}
}

func (w *Weather) Handle(ctx context.Context, e *model.Event) error {
	err := e.HandleTypeMessage(ctx, func(context.Context, *model.Event) error {
		if e.FilterText(triggerWeather) {
			return w.handleWeather(ctx, e)
		}

		return nil
	})
	if err != nil {
		return xerrors.Errorf("failed to handle type message: %w", err)
	}

	return nil
}

func (w *Weather) handleWeather(ctx context.Context, e *model.Event) error {
	imageURL, err := w.weather.ImageURL(ctx)
	if err != nil {
		return xerrors.Errorf("weather.Fetch: %w", err)
	}

	msg := w.message.Image(imageURL, imageURL)
	if err := w.bot.ReplyMessage(ctx, e, msg); err != nil {
		return xerrors.Errorf("bot.ReplyMessage: %w", err)
	}

	return nil
}
