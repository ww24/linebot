package accesslog

import (
	"bytes"
	"context"
	"log/slog"
	"time"

	"cloud.google.com/go/pubsub"

	"github.com/ww24/linebot/internal/accesslog/avro"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/log"
)

const resultChannelSize = 100

type Publisher interface {
	Publish(context.Context, *avro.AccessLog)
}

type NoopPublisher struct{}

func (*NoopPublisher) Publish(context.Context, *avro.AccessLog) {}

type PubSubPublisher struct {
	topic   *pubsub.Topic
	results chan *pubsub.PublishResult
}

func NewPublisher(p *pubsub.Client, cfg *config.AccessLog) (Publisher, func()) {
	if cfg.Topic == "" {
		return new(NoopPublisher), func() {} // noop
	}

	topic := p.Topic(cfg.Topic)
	topic.PublishSettings.DelayThreshold = time.Second
	publisher := &PubSubPublisher{
		topic:   topic,
		results: make(chan *pubsub.PublishResult, resultChannelSize),
	}
	ctx, cancel := context.WithCancel(context.Background())
	go publisher.worker(ctx)
	stop := func() {
		topic.Stop()
		cancel()
	}
	return publisher, stop
}

func (p *PubSubPublisher) Publish(ctx context.Context, al *avro.AccessLog) {
	buf := new(bytes.Buffer)
	if err := al.Serialize(buf); err != nil {
		slog.ErrorContext(ctx, "accesslog: failed to serialize access log", log.Err(err))
		return
	}
	result := p.topic.Publish(ctx, &pubsub.Message{
		Data: buf.Bytes(),
	})
	select {
	case p.results <- result:
	default:
		slog.InfoContext(ctx, "accesslog: publish results is sampled")
	}
}

func (p *PubSubPublisher) worker(ctx context.Context) {
	for result := range p.results {
		if _, err := result.Get(ctx); err != nil {
			slog.ErrorContext(ctx, "accesslog: failed to publish access log", log.Err(err))
		}
	}
}
