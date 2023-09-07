package main

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/tracer"
)

//nolint:gochecknoglobals
var tc = tracer.NewConfig(serviceName, version)

type bot struct {
	conf    *config.LINEBot
	handler http.Handler
}

func newBot(
	conf *config.LINEBot,
	handler http.Handler,
	_ trace.TracerProvider,
) *bot {
	return &bot{
		conf:    conf,
		handler: handler,
	}
}
