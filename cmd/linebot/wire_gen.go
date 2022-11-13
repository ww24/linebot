// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/infra/external/linebot"
	"github.com/ww24/linebot/infra/firestore"
	"github.com/ww24/linebot/infra/gcs"
	"github.com/ww24/linebot/infra/scheduler"
	"github.com/ww24/linebot/interactor"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/nl"
	"github.com/ww24/linebot/presentation/http"
	"github.com/ww24/linebot/tracer"
)

// Injectors from wire.go:

func register(contextContext context.Context) (*bot, func(), error) {
	lineBot, err := config.NewLINEBot()
	if err != nil {
		return nil, nil, err
	}
	logger, err := newLogger(contextContext)
	if err != nil {
		return nil, nil, err
	}
	linebotLINEBot, err := linebot.NewLINEBot(lineBot)
	if err != nil {
		return nil, nil, err
	}
	messageProviderSet := linebot.NewMessageProviderSet()
	botImpl := service.NewBot(linebotLINEBot, messageProviderSet)
	authorizer, err := http.NewAuthorizer(contextContext, lineBot)
	if err != nil {
		return nil, nil, err
	}
	tracerConfig := _wireConfigValue
	otel, err := config.NewOtel()
	if err != nil {
		return nil, nil, err
	}
	spanExporter := tracer.NewCloudTraceExporter()
	tracerProvider, cleanup := tracer.New(tracerConfig, otel, spanExporter)
	client, err := firestore.New(contextContext, tracerProvider)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	conversation := firestore.NewConversation(client)
	conversationImpl := service.NewConversation(conversation)
	shopping := firestore.NewShopping(conversation)
	shoppingImpl := service.NewShopping(conversation, shopping)
	parser, err := nl.NewParser()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	interactorShopping := interactor.NewShopping(conversationImpl, shoppingImpl, parser, messageProviderSet, botImpl)
	reminder := firestore.NewReminder(conversation)
	schedulerScheduler, err := scheduler.New(contextContext, lineBot)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	reminderImpl := service.NewReminder(reminder, schedulerScheduler)
	time, err := config.NewTime()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	interactorReminder := interactor.NewReminder(conversationImpl, reminderImpl, messageProviderSet, botImpl, time)
	gcsClient, err := gcs.New(contextContext)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	storage, err := config.NewStorage()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	weatherImageStore, err := gcs.NewWeatherImageStore(gcsClient, storage, time)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	weatherImpl := service.NewWeather(weatherImageStore, time)
	weather, err := interactor.NewWeather(weatherImpl, messageProviderSet, botImpl, lineBot)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	eventHandler, err := interactor.NewEventHandler(interactorShopping, interactorReminder, weather, conversationImpl, reminderImpl, messageProviderSet, botImpl, lineBot)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	imageStore, err := gcs.NewImageStore(gcsClient, storage)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	image := interactor.NewImage(imageStore)
	handler, err := http.NewHandler(logger, botImpl, authorizer, eventHandler, image)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	mainBot := newBot(lineBot, handler, tracerProvider)
	return mainBot, func() {
		cleanup()
	}, nil
}

var (
	_wireConfigValue = tc
)
