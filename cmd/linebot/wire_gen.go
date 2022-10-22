// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/ww24/linebot/config"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/infra/external/linebot"
	"github.com/ww24/linebot/infra/external/weather"
	"github.com/ww24/linebot/infra/firestore"
	"github.com/ww24/linebot/infra/gcs"
	"github.com/ww24/linebot/infra/scheduler"
	"github.com/ww24/linebot/interactor"
	"github.com/ww24/linebot/nl"
	"github.com/ww24/linebot/presentation/http"
	"github.com/ww24/linebot/tracer"
)

// Injectors from wire.go:

func register(contextContext context.Context) (*bot, func(), error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}
	logger, err := newLogger(contextContext)
	if err != nil {
		return nil, nil, err
	}
	lineBot, err := linebot.NewLINEBot(configConfig)
	if err != nil {
		return nil, nil, err
	}
	messageProviderSet := linebot.NewMessageProviderSet()
	botImpl := service.NewBot(lineBot, messageProviderSet)
	authorizer, err := http.NewAuthorizer(contextContext, configConfig)
	if err != nil {
		return nil, nil, err
	}
	client, err := firestore.New(contextContext)
	if err != nil {
		return nil, nil, err
	}
	conversation := firestore.NewConversation(client)
	conversationImpl := service.NewConversation(conversation)
	shopping := firestore.NewShopping(conversation)
	shoppingImpl := service.NewShopping(conversation, shopping)
	parser, err := nl.NewParser()
	if err != nil {
		return nil, nil, err
	}
	interactorShopping := interactor.NewShopping(conversationImpl, shoppingImpl, parser, messageProviderSet, botImpl)
	reminder := firestore.NewReminder(conversation)
	schedulerScheduler, err := scheduler.New(contextContext, configConfig)
	if err != nil {
		return nil, nil, err
	}
	reminderImpl := service.NewReminder(reminder, schedulerScheduler)
	interactorReminder := interactor.NewReminder(conversationImpl, reminderImpl, messageProviderSet, botImpl, configConfig)
	weatherWeather, err := weather.NewWeather(configConfig)
	if err != nil {
		return nil, nil, err
	}
	gcsClient, err := gcs.New(contextContext)
	if err != nil {
		return nil, nil, err
	}
	weatherImageStore, err := gcs.NewWeatherImageStore(gcsClient, configConfig)
	if err != nil {
		return nil, nil, err
	}
	weatherImpl := service.NewWeather(weatherWeather, weatherImageStore, configConfig)
	interactorWeather := interactor.NewWeather(weatherImpl, messageProviderSet, botImpl)
	eventHandler, err := interactor.NewEventHandler(interactorShopping, interactorReminder, interactorWeather, conversationImpl, reminderImpl, messageProviderSet, botImpl, configConfig)
	if err != nil {
		return nil, nil, err
	}
	imageStore, err := gcs.NewImageStore(gcsClient, configConfig)
	if err != nil {
		return nil, nil, err
	}
	image := interactor.NewImage(imageStore)
	handler, err := http.NewHandler(logger, botImpl, authorizer, eventHandler, image)
	if err != nil {
		return nil, nil, err
	}
	tracerConfig := _wireConfigValue
	spanExporter, err := tracer.NewCloudTraceExporter()
	if err != nil {
		return nil, nil, err
	}
	tracerProvider, cleanup := tracer.New(tracerConfig, configConfig, spanExporter)
	mainBot := newBot(configConfig, handler, tracerProvider)
	return mainBot, func() {
		cleanup()
	}, nil
}

var (
	_wireConfigValue = tc
)
