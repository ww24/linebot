package main

import (
	"context"
	"net/http"

	"github.com/ww24/linebot/logger"
)

type bot struct {
	log     *logger.Logger
	handler http.Handler
}

func newBot(
	log *logger.Logger,
	handler http.Handler,
) *bot {
	return &bot{
		log:     log,
		handler: handler,
	}
}

func newLogger(ctx context.Context) (*logger.Logger, error) {
	return logger.New(ctx, serviceName, version)
}
