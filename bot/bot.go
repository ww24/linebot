package bot

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/google/wire"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/ww24/linebot/domain/model"
)

var Set = wire.NewSet(
	New,
	NewShoppingService,
)

type Handler interface {
	Handle(context.Context, *Bot, *linebot.Event) error
}

type Bot struct {
	cli                    *linebot.Client
	allowedConversationIDs map[model.ConversationID]struct{}
	handlers               []Handler
}

type Config struct {
	ChannelSecret   string
	ChannelToken    string
	ConversationIDs []string
}

func New(
	c Config,
	shopping *ShoppingService,
) (*Bot, error) {
	hc := &http.Client{}
	cli, err := linebot.New(c.ChannelSecret, c.ChannelToken, linebot.WithHTTPClient(hc))
	if err != nil {
		return nil, err
	}

	allowedConversationIDs := make(map[model.ConversationID]struct{}, len(c.ConversationIDs))
	for _, id := range c.ConversationIDs {
		allowedConversationIDs[model.ConversationID(id)] = struct{}{}
	}

	bot := &Bot{
		cli:                    cli,
		allowedConversationIDs: allowedConversationIDs,
		handlers: []Handler{
			shopping,
		},
	}
	return bot, nil
}

func (b *Bot) HandleRequest(r *http.Request) error {
	events, err := b.cli.ParseRequest(r)
	if err != nil {
		return err
	}

	ctx := r.Context()

	for _, e := range events {
		if !b.filter(e) {
			log.Printf("ConversationID: %s\n", ConversationID(e.Source))
			return nil
		}

		for _, handler := range b.handlers {
			if err := handler.Handle(ctx, b, e); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bot) replyTestMessage(ctx context.Context, e *linebot.Event, str string) error {
	msg := linebot.NewTextMessage(str)
	c := b.cli.ReplyMessage(e.ReplyToken, msg)
	if _, err := c.WithContext(ctx).Do(); err != nil {
		return err
	}
	return nil
}

func (b *Bot) filter(e *linebot.Event) bool {
	cID := ConversationID(e.Source)
	if len(b.allowedConversationIDs) > 0 {
		if _, ok := b.allowedConversationIDs[cID]; !ok {
			return false
		}
	}

	return true
}

func (b *Bot) filterText(e *linebot.Event, target string) bool {
	text, ok := e.Message.(*linebot.TextMessage)
	if ok && strings.Contains(text.Text, target) {
		return true
	}

	return false
}

func (b *Bot) readTextLines(e *linebot.Event) []string {
	text, ok := e.Message.(*linebot.TextMessage)
	if !ok {
		return nil
	}

	lines := strings.Split(text.Text, "\n")
	ret := make([]string, 0, len(lines))
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			ret = append(ret, line)
		}
	}

	return ret
}
