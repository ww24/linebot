package config

import (
	"os"
	"strings"

	"github.com/google/wire"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewConfig,
	wire.Bind(new(repository.Config), new(*Config)),
)

// Config implements repository.Config.
type Config struct {
	lineChannelSecret string
	lineChannelToken  string
	conversationIDs   *ConversationIDs
	addr              string
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

	return &Config{
		lineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		lineChannelToken:  os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		conversationIDs:   conversationIDs,
		addr:              addr,
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
