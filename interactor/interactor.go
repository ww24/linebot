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
	NewReminder,
	NewShopping,
)

type EventHandler struct {
	handlers         []repository.Handler
	scheduleHandlers []repository.ScheduleHandler
	remindHandlers   []repository.RemindHandler
	reminder         service.Reminder
	conf             repository.Config
	bot              service.Bot
	message          repository.MessageProviderSet
}

func NewEventHandler(
	shoppingInteractor *Shopping,
	reminderInteractor *Reminder,
	reminder service.Reminder,
	message repository.MessageProviderSet,
	bot service.Bot,
	conf repository.Config,
) (*EventHandler, error) {
	return &EventHandler{
		handlers: []repository.Handler{
			shoppingInteractor,
			reminderInteractor,
		},
		scheduleHandlers: []repository.ScheduleHandler{
			reminderInteractor,
		},
		remindHandlers: []repository.RemindHandler{
			shoppingInteractor,
		},
		reminder: reminder,
		conf:     conf,
		bot:      bot,
		message:  message,
	}, nil
}

func (h *EventHandler) Handle(ctx context.Context, events []*model.Event) error {
	dl := logger.DefaultLogger(ctx)

	for _, e := range events {
		if !h.conf.ConversationIDs().Available(e.ConversationID()) {
			dl.Warn("not allowed conversation",
				zap.String("ConversationID", e.ConversationID().String()),
			)
			return nil
		}

		for _, handler := range h.handlers {
			if err := handler.Handle(ctx, e); err != nil {
				dl.Error("failed to handle event", zap.Error(err))
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
	dl := logger.DefaultLogger(ctx)

	for _, handler := range h.scheduleHandlers {
		if err := handler.HandleSchedule(ctx); err != nil {
			dl.Error("failed to handle schedule", zap.Error(err))
			return xerrors.Errorf("failed to handle schedule: %w", err)
		}
	}

	return nil
}

func (h *EventHandler) HandleReminder(ctx context.Context, itemIDJSON *model.ReminderItemIDJSON) error {
	dl := logger.DefaultLogger(ctx)

	item, err := h.reminder.Get(ctx, model.ConversationID(itemIDJSON.ConversationID), model.ReminderItemID(itemIDJSON.ItemID))
	if err != nil {
		dl.Error("failed to get reminder item", zap.Error(err))
		return xerrors.Errorf("failed to get reminder item: %w", err)
	}

	for _, handler := range h.remindHandlers {
		if err := handler.HandleReminder(ctx, item); err != nil {
			dl.Error("failed to handle reminder", zap.Error(err))
			return xerrors.Errorf("failed to handle reminder: %w", err)
		}
	}

	return nil
}
