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
	message  repository.MessageProviderSet
}

func NewEventHandler(
	conversation service.Conversation,
	shopping service.Shopping,
	nlParser repository.NLParser,
	message repository.MessageProviderSet,
	bot service.Bot,
	conf repository.Config,
	log *logger.Logger,
) (*EventHandler, error) {
	return &EventHandler{
		handlers: []repository.Handler{
			NewShopping(conversation, shopping, nlParser, message, bot),
		},
		conf:    conf,
		log:     log,
		bot:     bot,
		message: message,
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
				cl.Error("failed to handle event", zap.Error(err))
				return h.handleError(ctx, e)
			}
		}
	}

	return nil
}

func (h *EventHandler) handleError(ctx context.Context, e *model.Event) error {
	msg := h.message.Text("予期せぬエラーが発生しました")
	if err := h.bot.PushMessage(ctx, e.ConversationID(), msg); err != nil {
		return xerrors.Errorf("failed to reply text message: %w", err)
	}

	return nil
}
