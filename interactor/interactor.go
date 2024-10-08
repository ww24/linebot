package interactor

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/wire"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/internal/code"
	"github.com/ww24/linebot/internal/config"
	"github.com/ww24/linebot/log"
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
	conversation     service.Conversation
	reminder         service.Reminder
	conversationIDs  *config.ConversationIDs
	bot              service.Bot
	message          repository.MessageProviderSet
}

func NewEventHandler(
	shoppingInteractor *Shopping,
	reminderInteractor *Reminder,
	weatherInteractor *Weather,
	conversation service.Conversation,
	reminder service.Reminder,
	message repository.MessageProviderSet,
	bot service.Bot,
	conf *config.LINEBot,
) (*EventHandler, error) {
	return &EventHandler{
		handlers: []repository.Handler{
			shoppingInteractor,
			reminderInteractor,
			weatherInteractor,
		},
		scheduleHandlers: []repository.ScheduleHandler{
			reminderInteractor,
		},
		remindHandlers: []repository.RemindHandler{
			shoppingInteractor,
		},
		conversation:    conversation,
		reminder:        reminder,
		conversationIDs: conf.ConversationIDs(),
		bot:             bot,
		message:         message,
	}, nil
}

func (h *EventHandler) Handle(ctx context.Context, events []*model.Event) error {
	for _, e := range events {
		if !h.conversationIDs.Available(e.ConversationID()) {
			slog.WarnContext(ctx, "interactor: not allowed conversation",
				slog.String("ConversationID", e.ConversationID().String()),
			)
			return nil
		}

		status, err := h.conversation.GetStatus(ctx, e.ConversationID())
		if err != nil {
			return xerrors.Errorf("failed to get status: %w", err)
		}
		e.SetStatus(status.Type)

		for _, handler := range h.handlers {
			if err := handler.Handle(ctx, e); err != nil {
				if errors.Is(err, errResponseReturned) {
					return nil
				}

				slog.ErrorContext(ctx, "interactor: failed to handle event", log.Err(err))
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
	item, err := h.reminder.Get(ctx, model.ConversationID(itemIDJSON.ConversationID), model.ReminderItemID(itemIDJSON.ItemID))
	if err != nil {
		if code.From(err) == code.NotFound {
			slog.InfoContext(ctx, "interactor: reminder item not found",
				slog.String("ConversationID", itemIDJSON.ConversationID),
				slog.String("ItemID", itemIDJSON.ItemID),
				log.Err(err),
			)
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
