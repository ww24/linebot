//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOPACKAGE/mock_$GOFILE -package=mock_repository

package repository

import (
	"net/url"
	"time"

	"github.com/ww24/linebot/domain/model"
)

type Config interface {
	LINEChannelSecret() string
	LINEChannelToken() string
	ConversationIDs() ConversationIDs
	Addr() string
	CloudTasksLocation() string
	CloudTasksQueue() string
	ServiceEndpoint(path string) (*url.URL, error)
	WeatherAPI() string
	ImageBucket() string
	DefaultLocation() *time.Location
}

type ConversationIDs interface {
	List() []model.ConversationID
	Available(model.ConversationID) bool
}
