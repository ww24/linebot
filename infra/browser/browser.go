package browser

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/wire"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/logger"
)

const (
	windowWidth  = 1280
	windowHeight = 960
	dialTimeout  = 10 * time.Second
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewBrowser,
	wire.Bind(new(repository.Browser), new(*Browser)),
)

type Browser struct {
	timeout time.Duration
}

func NewBrowser(conf repository.Config) *Browser {
	return &Browser{
		timeout: conf.BrowserTimeout(),
	}
}

func (b *Browser) Screenshot(ctx context.Context, target *url.URL, targetSelector string) (io.Reader, int, error) {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(windowWidth, windowHeight),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()
	taskCtx, cancel := chromedp.NewContext(allocCtx,
		chromedp.WithBrowserOption(
			chromedp.WithDialTimeout(dialTimeout),
		),
	)
	defer cancel()

	dl := logger.DefaultLogger(ctx)
	dl.Info("capture screenshot",
		zap.String("target", target.String()),
		zap.String("selector", targetSelector),
	)

	var buf []byte
	tasks := chromedp.Tasks{
		chromedp.Navigate(target.String()),
		chromedp.Screenshot(targetSelector, &buf, chromedp.ByID),
	}

	if err := chromedp.Run(taskCtx, tasks...); err != nil {
		return nil, 0, xerrors.Errorf("chromedp.Run: %w", err)
	}

	return bytes.NewReader(buf), len(buf), nil
}
