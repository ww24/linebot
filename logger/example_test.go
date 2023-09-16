package logger_test

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/ww24/linebot/logger"
)

func Example() {
	ctx := context.Background()
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	// initialize root logger
	if err := logger.SetConfigWithWriter("service-name", "v1.0.0", os.Stdout); err != nil {
		panic(err)
	}

	dl := logger.Default(ctx)
	dl = dl.WithLogger(dl.WithOptions(zap.WithClock(staticClock(testTime))))

	// info log
	dl.Info("message", zap.String("key", "value"))

	// Output:
	// {"severity":"INFO","timestamp":"2023-01-01T00:00:00Z","caller":"logger/example_test.go:26","message":"message","serviceContext":{"service":"service-name","version":"v1.0.0"},"key":"value","logging.googleapis.com/labels":{},"logging.googleapis.com/sourceLocation":{"file":"github.com/ww24/linebot/logger/example_test.go","line":"26","function":"github.com/ww24/linebot/logger_test.Example"}}
}

type staticClock time.Time

func (c staticClock) Now() time.Time {
	return time.Time(c)
}

func (c staticClock) NewTicker(time.Duration) *time.Ticker {
	tk := time.NewTicker(time.Second)
	tk.Stop()
	return tk
}
