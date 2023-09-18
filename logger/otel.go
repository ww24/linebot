package logger

import (
	"os"

	"github.com/go-logr/zapr"
	"go.opentelemetry.io/otel"
)

func init() {
	logLevel := getLogLevel("OTEL_LOG_SEVERITY_LEVEL")
	logger, err := newLogger(os.Stderr, logLevel)
	if err != nil {
		panic(err)
	}
	otel.SetLogger(zapr.NewLogger(logger.Logger))
}
