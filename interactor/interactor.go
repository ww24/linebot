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
	handlers         []repository.Handler
	scheduleHandlers []repository.ScheduleHandler
	remindHandlers   []repository.RemindHandler
	reminder         service.Reminder
	conf             repository.Config
	log              *logger.Logger
	bot              service.Bot
	message          repository.MessageProviderSet
}

func NewEventHandler(
	conversation service.Conversation,
	shopping service.Shopping,
	reminder service.Reminder,
	nlParser repository.NLParser,
	message repository.MessageProviderSet,
	bot service.Bot,
	conf repository.Config,
	log *logger.Logger,
) (*EventHandler, error) {
	hReminder := NewReminder(conversation, reminder, message, bot)
	hShopping := NewShopping(conversation, shopping, nlParser, message, bot)
	return &EventHandler{
		handlers: []repository.Handler{
			hShopping,
			hReminder,
		},
		scheduleHandlers: []repository.ScheduleHandler{
			hReminder,
		},
		remindHandlers: []repository.RemindHandler{
			hShopping,
		},
		reminder: reminder,
		conf:     conf,
		log:      log,
		bot:      bot,
		message:  message,
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

func (h *EventHandler) HandleSchedule(ctx context.Context) error {
	cl := h.log.WithTraceFromContext(ctx)

	for _, handler := range h.scheduleHandlers {
		if err := handler.HandleSchedule(ctx); err != nil {
			cl.Error("failed to handle schedule", zap.Error(err))
			return xerrors.Errorf("failed to handle schedule: %w", err)
		}
	}

	return nil
}

func (h *EventHandler) HandleReminder(ctx context.Context, itemIDJSON *model.ReminderItemIDJSON) error {
	cl := h.log.WithTraceFromContext(ctx)

	item, err := h.reminder.Get(ctx, model.ConversationID(itemIDJSON.ConversationID), model.ReminderItemID(itemIDJSON.ItemID))
	if err != nil {
		cl.Error("failed to get reminder item", zap.Error(err))
		return xerrors.Errorf("failed to get reminder item: %w", err)
	}

	for _, handler := range h.remindHandlers {
		if err := handler.HandleReminder(ctx, item); err != nil {
			cl.Error("failed to handle reminder", zap.Error(err))
			return xerrors.Errorf("failed to handle reminder: %w", err)
		}
	}

	return nil
}
