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
)

// Injectors from wire.go:

func register(contextContext context.Context) (*bot, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	logger, err := newLogger(contextContext)
	if err != nil {
		return nil, err
	}
	lineBot, err := linebot.NewLINEBot(configConfig)
	if err != nil {
		return nil, err
	}
	messageProviderSet := linebot.NewMessageProviderSet()
	botImpl := service.NewBot(lineBot, messageProviderSet)
	client, err := firestore.New(contextContext)
	if err != nil {
		return nil, err
	}
	conversation := firestore.NewConversation(client)
	conversationImpl := service.NewConversation(conversation)
	shoppingImpl := service.NewShopping(conversation)
	parser, err := nl.NewParser()
	if err != nil {
		return nil, err
	}
	shopping := interactor.NewShopping(conversationImpl, shoppingImpl, parser, messageProviderSet, botImpl)
	reminder := firestore.NewReminder(conversation)
	schedulerScheduler, err := scheduler.New(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	reminderImpl := service.NewReminder(reminder, schedulerScheduler)
	interactorReminder := interactor.NewReminder(conversationImpl, reminderImpl, messageProviderSet, botImpl, configConfig)
	weatherWeather, err := weather.NewWeather(configConfig)
	if err != nil {
		return nil, err
	}
	gcsClient, err := gcs.New(contextContext)
	if err != nil {
		return nil, err
	}
	weatherImageStore, err := gcs.NewWeatherImageStore(gcsClient, configConfig)
	if err != nil {
		return nil, err
	}
	weatherImpl := service.NewWeather(weatherWeather, weatherImageStore)
	interactorWeather := interactor.NewWeather(weatherImpl, messageProviderSet, botImpl)
	eventHandler, err := interactor.NewEventHandler(shopping, interactorReminder, interactorWeather, reminderImpl, messageProviderSet, botImpl, configConfig)
	if err != nil {
		return nil, err
	}
	imageStore, err := gcs.NewImageStore(gcsClient, configConfig)
	if err != nil {
		return nil, err
	}
	image := interactor.NewImage(imageStore)
	handler := http.NewHandler(logger, botImpl, eventHandler, image)
	mainBot := newBot(configConfig, handler)
	return mainBot, nil
}
