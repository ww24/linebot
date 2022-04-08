package config

import (
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/wire"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
)

const (
	defaultTimezoneOffset = 9 * 60 * 60
)

var (
	//nolint: gochecknoglobals
	defaultLocation = time.FixedZone("Asia/Tokyo", defaultTimezoneOffset)
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewConfig,
	wire.Bind(new(repository.Config), new(*Config)),
)

// Config implements repository.Config.
type Config struct {
	lineChannelSecret  string
	lineChannelToken   string
	conversationIDs    *ConversationIDs
	addr               string
	cloudTasksLocation string
	cloudTasksQueue    string
	serviceEndpoint    *url.URL
	weatherAPI         string
	imageBucket        string
	defaultTimezone    string
}

func NewConfig() *Config {
	var conversationIDs = &ConversationIDs{
		list: make([]model.ConversationID, 0),
		set:  make(map[model.ConversationID]struct{}),
	}
	for _, id := range strings.Split(os.Getenv("ALLOW_CONV_IDS"), ",") {
		conversationID := model.ConversationID(id)
		conversationIDs.list = append(conversationIDs.list, conversationID)
		conversationIDs.set[conversationID] = struct{}{}
	}

	addr := ":8000"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	var serviceEndpoint *url.URL
	if u, err := url.Parse(os.Getenv("SERVICE_ENDPOINT")); err == nil {
		serviceEndpoint = u
	}

	return &Config{
		lineChannelSecret:  os.Getenv("LINE_CHANNEL_SECRET"),
		lineChannelToken:   os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		conversationIDs:    conversationIDs,
		addr:               addr,
		cloudTasksLocation: os.Getenv("CLOUD_TASKS_LOCATION"),
		cloudTasksQueue:    os.Getenv("CLOUD_TASKS_QUEUE"),
		serviceEndpoint:    serviceEndpoint,
		weatherAPI:         os.Getenv("WEATHER_API"),
		imageBucket:        os.Getenv("IMAGE_BUCKET"),
		defaultTimezone:    os.Getenv("DEFAULT_TIMEZONE"),
	}
}

func (c *Config) LINEChannelSecret() string {
	return c.lineChannelSecret
}

func (c *Config) LINEChannelToken() string {
	return c.lineChannelToken
}

func (c *Config) ConversationIDs() repository.ConversationIDs {
	return c.conversationIDs
}

// ConversationIDs implements repository.ConversationIDs.
type ConversationIDs struct {
	list []model.ConversationID
	set  map[model.ConversationID]struct{}
}

func (c *ConversationIDs) List() []model.ConversationID {
	return c.list
}

func (c *ConversationIDs) Available(conversationID model.ConversationID) bool {
	// return true if conversationIDs is empty
	if len(c.set) == 0 {
		return true
	}

	_, ok := c.set[conversationID]
	return ok
}

func (c *Config) Addr() string {
	return c.addr
}

func (c *Config) CloudTasksLocation() string {
	return c.cloudTasksLocation
}

func (c *Config) CloudTasksQueue() string {
	return c.cloudTasksQueue
}

func (c *Config) ServiceEndpoint(path string) (*url.URL, error) {
	if c.serviceEndpoint == nil {
		return nil, xerrors.New("service endpoint is not configured")
	}
	r, err := url.Parse(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse path: %w", err)
	}
	return c.serviceEndpoint.ResolveReference(r), nil
}

func (c *Config) WeatherAPI() string {
	return c.weatherAPI
}

func (c *Config) ImageBucket() string {
	return c.imageBucket
}

func (c *Config) DefaultLocation() *time.Location {
	if c.defaultTimezone == "" {
		return defaultLocation
	}
	loc, err := time.LoadLocation(c.defaultTimezone)
	if err != nil {
		return defaultLocation
	}
	return loc
}
