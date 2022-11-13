package linebot

import (
	"context"
	"net/http"

	"github.com/google/wire"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/tracer"
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewLINEBot,
	wire.Bind(new(repository.Bot), new(*LINEBot)),
	NewMessageProviderSet,
	wire.Bind(new(repository.MessageProviderSet), new(*MessageProviderSet)),
)

// LINEBot implements repository.Bot.
type LINEBot struct {
	cli                    *linebot.Client
	allowedConversationIDs *config.ConversationIDs
}

func NewLINEBot(conf *config.LINEBot) (*LINEBot, error) {
	transport := tracer.HTTPTransport(http.DefaultTransport)
	hc := &http.Client{Transport: transport}
	cli, err := linebot.New(
		conf.LINEChannelSecret,
		conf.LINEChannelAccessToken,
		linebot.WithHTTPClient(hc),
	)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize LINE Bot client: %w", err)
	}

	return &LINEBot{
		cli:                    cli,
		allowedConversationIDs: conf.ConversationIDs(),
	}, nil
}

func (b *LINEBot) EventsFromRequest(r *http.Request) ([]*model.Event, error) {
	events, err := b.cli.ParseRequest(r)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse request: %w", err)
	}

	es := make([]*model.Event, 0, len(events))
	for _, event := range events {
		e := new(model.Event)
		e.Event = event
		e.SetStatus(model.ConversationStatusTypeNeutral)
		es = append(es, e)
	}

	return es, nil
}

func (b *LINEBot) ReplyMessage(ctx context.Context, e *model.Event, p repository.MessageProvider) error {
	msg := p.ToMessage()

	c := b.cli.ReplyMessage(e.ReplyToken, msg)
	if _, err := c.WithContext(ctx).Do(); err != nil {
		return xerrors.Errorf("failed to reply message: %w", err)
	}

	return nil
}

func (b *LINEBot) PushMessage(ctx context.Context, to model.ConversationID, p repository.MessageProvider) error {
	msg := p.ToMessage()

	c := b.cli.PushMessage(to.SourceID(), msg)
	if _, err := c.WithContext(ctx).Do(); err != nil {
		return xerrors.Errorf("failed to reply message: %w", err)
	}

	return nil
}
