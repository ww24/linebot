package main

import (
	"context"
	"net/http"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/logger"
)

type bot struct {
	config  repository.Config
	log     *logger.Logger
	handler http.Handler
}

func newBot(
	config repository.Config,
	log *logger.Logger,
	handler http.Handler,
) *bot {
	return &bot{
		config:  config,
		log:     log,
		handler: handler,
	}
}

func newLogger(ctx context.Context) (*logger.Logger, error) {
	return logger.New(ctx, serviceName, version)
}
