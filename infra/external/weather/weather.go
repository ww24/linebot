package weather

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/wire"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/xerrors"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"

	"github.com/ww24/linebot/domain/repository"
)

const (
	timeout = 10 * time.Second
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewWeather,
	wire.Bind(new(repository.Weather), new(*Weather)),
)

type Weather struct {
	endpoint string
	audience string
}

func NewWeather(conf repository.Config) (*Weather, error) {
	u, err := url.Parse(conf.WeatherAPI())
	if err != nil {
		return nil, xerrors.Errorf("failed to parse weather api url: %w", err)
	}
	audience := u.Scheme + "://" + u.Hostname() + "/"

	return &Weather{
		endpoint: conf.WeatherAPI(),
		audience: audience,
	}, nil
}

func (w *Weather) newTransport(ctx context.Context) (http.RoundTripper, error) {
	ts, err := idtoken.NewTokenSource(ctx, w.audience)
	if err != nil {
		return nil, xerrors.Errorf("failed to create token source: %w", err)
	}

	t, err := htransport.NewTransport(ctx, otelhttp.NewTransport(http.DefaultTransport), option.WithTokenSource(ts))
	if err != nil {
		return nil, xerrors.Errorf("failed to create idtoken client: %w", err)
	}

	return t, nil
}

// Fetch an weather map image.
func (w *Weather) Fetch(ctx context.Context) (io.ReadCloser, error) {
	t, err := w.newTransport(ctx)
	if err != nil {
		return nil, err
	}
	cli := &http.Client{
		Timeout:   timeout,
		Transport: t,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.endpoint, http.NoBody)
	if err != nil {
		return nil, xerrors.Errorf("failed to create http request: %w", err)
	}

	res, err := cli.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to get weather: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		return nil, xerrors.Errorf("weather api response is not ok: %s", res.Status)
	}

	return res.Body, nil
}
