package interactor

import (
	"context"

	"github.com/google/wire"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

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
	bot      service.Bot
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
		bot:  bot,
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

		for _, handler := range h.handlers {
			if err := handler.Handle(ctx, e); err != nil {
				return h.handleError(ctx, e)
			}
		}
	}

	return nil
}

func (h *EventHandler) handleError(ctx context.Context, e *model.Event) error {
	if err := h.bot.ReplyTextMessage(ctx, e, "予期せぬエラーが発生しました"); err != nil {
		return xerrors.Errorf("failed to reply text message: %w", err)
	}

	return nil
}
