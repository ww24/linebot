package interactor

import (
	"context"
	"errors"

	"github.com/google/wire"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/internal/code"
	"github.com/ww24/linebot/logger"
	"github.com/ww24/linebot/usecase"
)

var (
	errResponseReturned = xerrors.New("response returned")
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewEventHandler,
	wire.Bind(new(usecase.EventHandler), new(*EventHandler)),
	NewReminder,
	NewShopping,
	NewScreenshot,
	wire.Bind(new(usecase.ScreenshotHandler), new(*Screenshot)),
	NewWeather,
	NewImage,
	wire.Bind(new(usecase.ImageHandler), new(*Image)),
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
	weatherInteractor *Weather,
	reminder service.Reminder,
	message repository.MessageProviderSet,
	bot service.Bot,
	conf repository.Config,
) (*EventHandler, error) {
	return &EventHandler{
		handlers: []repository.Handler{
			shoppingInteractor,
			reminderInteractor,
			weatherInteractor,
		},
		scheduleHandlers: []repository.ScheduleHandler{
			reminderInteractor,
			weatherInteractor,
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
				if errors.Is(err, errResponseReturned) {
					return nil
				}

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
		return xerrors.Errorf("bot.PushMessage: %w", err)
	}

	return nil
}

func (h *EventHandler) HandleSchedule(ctx context.Context) error {
	for _, handler := range h.scheduleHandlers {
		if err := handler.HandleSchedule(ctx); err != nil {
			return xerrors.Errorf("failed to handle schedule: %w", err)
		}
	}

	return nil
}

func (h *EventHandler) HandleReminder(ctx context.Context, itemIDJSON *model.ReminderItemIDJSON) error {
	dl := logger.DefaultLogger(ctx)

	item, err := h.reminder.Get(ctx, model.ConversationID(itemIDJSON.ConversationID), model.ReminderItemID(itemIDJSON.ItemID))
	if err != nil {
		if code.From(err) == code.NotFound {
			dl.Info("reminder item not found", zap.Error(err))
			return nil
		}
		return xerrors.Errorf("failed to get reminder item: %w", err)
	}

	for _, handler := range h.remindHandlers {
		if err := handler.HandleReminder(ctx, item); err != nil {
			return xerrors.Errorf("failed to handle reminder: %w", err)
		}
	}

	return nil
}
