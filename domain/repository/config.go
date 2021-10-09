package repository

import "github.com/ww24/linebot/domain/model"

type Config interface {
	LINEChannelSecret() string
	LINEChannelToken() string
	ConversationIDs() ConversationIDs
}

type ConversationIDs interface {
	List() []model.ConversationID
	Available(model.ConversationID) bool
}
