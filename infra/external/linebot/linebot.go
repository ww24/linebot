package linebot

import (
	"context"
	"net/http"

	"github.com/google/wire"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
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
	allowedConversationIDs repository.ConversationIDs
}

func NewLINEBot(conf repository.Config) (*LINEBot, error) {
	hc := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	cli, err := linebot.New(
		conf.LINEChannelSecret(),
		conf.LINEChannelToken(),
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
