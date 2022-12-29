package accesslog

import (
	"bytes"
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"

	"github.com/ww24/linebot/internal/accesslog/avro"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/logger"
)

type Publisher interface {
	Publish(context.Context, *avro.AccessLog)
}

type NoopPublisher struct{}

func (*NoopPublisher) Publish(context.Context, *avro.AccessLog) {}

type PubSubPublisher struct {
	topic *pubsub.Topic
}

func NewPublisher(p *pubsub.Client, cfg *config.AccessLog) (Publisher, func()) {
	if cfg.Topic == "" {
		return new(NoopPublisher), func() {} // noop
	}

	topic := p.Topic(cfg.Topic)
	topic.PublishSettings.DelayThreshold = time.Second
	return &PubSubPublisher{
		topic: topic,
	}, topic.Stop
}

func (p *PubSubPublisher) Publish(ctx context.Context, al *avro.AccessLog) {
	buf := new(bytes.Buffer)
	if err := al.Serialize(buf); err != nil {
		dl := logger.DefaultLogger(ctx)
		dl.Error("failed to serialize access log", zap.Error(err))
		return
	}
	p.topic.Publish(ctx, &pubsub.Message{
		Data: buf.Bytes(),
	})
}
