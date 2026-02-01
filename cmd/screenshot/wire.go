//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/infra/browser"
	"github.com/ww24/linebot/infra/gcs"
	"github.com/ww24/linebot/interactor"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/tracer"
)

func register(
	ctx context.Context,
) (*job, func(), error) {
	wire.Build(
		config.Set,
		gcs.Set,
		browser.Set,
		service.Set,
		interactor.Set,
		wire.Value(tc),
		tracer.NewCloudTraceExporter,
		tracer.New,
		newJob,
	)
	return nil, nil, nil
}
