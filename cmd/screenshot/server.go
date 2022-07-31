package main

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/logger"
)

const readHeaderTimeout = 10 * time.Second

type server struct {
	config repository.Config
	srv    *http.Server
}

func newServer(conf repository.Config, handler http.Handler) *server {
	return &server{
		config: conf,
		srv: &http.Server{
			Handler:           handler,
			Addr:              conf.Addr(),
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
}

func newLogger(ctx context.Context) (*logger.Logger, error) {
	l, err := logger.New(ctx, serviceName, version)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize logger: %w", err)
	}

	return l, nil
}
