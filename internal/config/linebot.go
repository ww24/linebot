package config

import (
	"strconv"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
)

type LINEBot struct {
	LINEChannelSecret          string   `split_words:"true" required:"true"`
	LINEChannelAccessToken     string   `split_words:"true" required:"true"`
	AllowConvIDs               []string `split_words:"true"`
	Port                       int      `split_words:"true" default:"8000"`
	CloudTasksLocation         string   `split_words:"true" required:"true"`
	CloudTasksQueue            string   `split_words:"true" required:"true"`
	InvokerServiceAccountID    string   `split_words:"true" required:"true"`
	InvokerServiceAccountEmail string   `split_words:"true" required:"true"`
}

func NewLINEBot() (*LINEBot, error) {
	var conf LINEBot
	if err := envconfig.Process("", &conf); err != nil {
		return nil, xerrors.Errorf("failed to parse linebot config: %w", err)
	}
	return &conf, nil
}

func (c *LINEBot) Addr() string {
	return ":" + strconv.Itoa(c.Port)
}

func (c *LINEBot) ConversationIDs() *ConversationIDs {
	conversationIDs := &ConversationIDs{
		list: make([]model.ConversationID, 0),
		set:  make(map[model.ConversationID]struct{}),
	}
	for _, id := range c.AllowConvIDs {
		conversationID := model.ConversationID(id)
		if _, ok := conversationIDs.set[conversationID]; ok {
			continue
		}
		conversationIDs.list = append(conversationIDs.list, conversationID)
		conversationIDs.set[conversationID] = struct{}{}
	}
	return conversationIDs
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
