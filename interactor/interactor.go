package interactor

import (
	"context"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/usecase"
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewEventHandler,
	wire.Bind(new(usecase.EventHandler), new(*EventHandler)),
)

type EventHandler struct {
	handlers []repository.Handler
	conf     repository.Config
	log      *logger.Logger
}

func NewEventHandler(
	shopping service.Shopping,
	nlParser repository.NLParser,
	message repository.MessageProviderSet,
	bot service.Bot,
	conf repository.Config,
	log *logger.Logger,
) (*EventHandler, error) {
	return &EventHandler{
		handlers: []repository.Handler{
			NewShopping(shopping, nlParser, message, bot),
		},
		conf: conf,
		log:  log,
	}, nil
}

func (h *EventHandler) Handle(ctx context.Context, events []*model.Event) error {
	cl := h.log.WithTraceFromContext(ctx)

	for _, e := range events {
		if !h.conf.ConversationIDs().Available(e.ConversationID()) {
			cl.Warn("not allowed conversation",
				zap.String("ConversationID", e.ConversationID().String()),
			)
			return nil
		}

		for _, h := range h.handlers {
			if err := h.Handle(ctx, e); err != nil {
				return err
			}
		}
	}

	return nil
}
