//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/infra/external/linebot"
	"github.com/ww24/linebot/infra/firestore"
	"github.com/ww24/linebot/infra/gcs"
	"github.com/ww24/linebot/infra/pubsub"
	"github.com/ww24/linebot/infra/scheduler"
	"github.com/ww24/linebot/interactor"
	"github.com/ww24/linebot/internal/accesslog"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/nl"
	"github.com/ww24/linebot/presentation/http"
	"github.com/ww24/linebot/tracer"
)

func register(
	context.Context,
) (*bot, func(), error) {
	wire.Build(
		config.Set,
		firestore.Set,
		scheduler.Set,
		linebot.Set,
		gcs.Set,
		service.Set,
		nl.Set,
		interactor.Set,
		pubsub.Set,
		accesslog.Set,
		http.Set,
		wire.Value(tc),
		tracer.NewCloudTraceExporter,
		tracer.New,
		newBot,
	)
	return nil, nil, nil
}
