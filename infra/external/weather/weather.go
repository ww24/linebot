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
	client   *http.Client
}

func NewWeather(conf repository.Config) (*Weather, error) {
	client := &http.Client{
		Timeout:   timeout,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	u, err := url.Parse(conf.WeatherAPI())
	if err != nil {
		return nil, xerrors.Errorf("failed to parse weather api url: %w", err)
	}
	audience := u.Scheme + "://" + u.Hostname() + "/"

	return &Weather{
		endpoint: conf.WeatherAPI(),
		audience: audience,
		client:   client,
	}, nil
}

// Fetch an weather map image.
func (w *Weather) Fetch(ctx context.Context) (io.ReadCloser, error) {
	cli, err := idtoken.NewClient(ctx, w.audience, idtoken.WithHTTPClient(w.client))
	if err != nil {
		return nil, xerrors.Errorf("failed to create idtoken client: %w", err)
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
