package service

import (
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/logger"
)

const (
	syncInterval = 2 * time.Hour
)

type Reminder interface {
	Add(context.Context, *model.ReminderItem) error
	List(context.Context, model.ConversationID) (model.ReminderItems, error)
	Get(context.Context, model.ConversationID, model.ReminderItemID) (*model.ReminderItem, error)
	Delete(context.Context, model.ConversationID, model.ReminderItemID) error
	ListAll(context.Context) (model.ReminderItems, error)
	SyncSchedule(context.Context, model.ReminderItems) error
}

type ReminderImpl struct {
	reminder  repository.Reminder
	scheduler repository.ScheduleSynchronizer
}

func NewReminder(
	reminder repository.Reminder,
	scheduler repository.ScheduleSynchronizer,
) *ReminderImpl {
	return &ReminderImpl{
		reminder:  reminder,
		scheduler: scheduler,
	}
}

func (r *ReminderImpl) Add(ctx context.Context, item *model.ReminderItem) error {
	ctx, span := tracer.Start(ctx, "Reminder#Add")
	defer span.End()

	if err := r.reminder.Add(ctx, item); err != nil {
		return xerrors.Errorf("failed to add a reminder item: %w", err)
	}

	now := time.Now()
	items := model.ReminderItems{item}.FilterNextSchedule(now, syncInterval)
	for _, item := range items {
		if err := r.scheduler.Create(ctx, item.ConversationID, item, now); err != nil {
			return xerrors.Errorf("failed to create a schedule: %w", err)
		}
	}

	return nil
}

func (r *ReminderImpl) List(ctx context.Context, conversationID model.ConversationID) (model.ReminderItems, error) {
	ctx, span := tracer.Start(ctx, "Reminder#List")
	defer span.End()

	items, err := r.reminder.List(ctx, conversationID)
	if err != nil {
		return nil, xerrors.Errorf("failed to list reminder items: %w", err)
	}
	return items, nil
}

func (r *ReminderImpl) Get(ctx context.Context, conversationID model.ConversationID, itemID model.ReminderItemID) (*model.ReminderItem, error) {
	ctx, span := tracer.Start(ctx, "Reminder#Get")
	defer span.End()

	item, err := r.reminder.Get(ctx, conversationID, itemID)
	if err != nil {
		return nil, xerrors.Errorf("failed to get a reminder item: %w", err)
	}
	return item, nil
}

func (r *ReminderImpl) Delete(ctx context.Context, conversationID model.ConversationID, itemID model.ReminderItemID) error {
	ctx, span := tracer.Start(ctx, "Reminder#Delete")
	defer span.End()

	item, err := r.reminder.Get(ctx, conversationID, itemID)
	if err != nil {
		return xerrors.Errorf("failed to get a reminder item: %w", err)
	}
	if err := r.reminder.Delete(ctx, conversationID, itemID); err != nil {
		return xerrors.Errorf("failed to delete a reminder item: %w", err)
	}

	now := time.Now()
	items := model.ReminderItems{item}.FilterNextSchedule(now, syncInterval)
	for _, item := range items {
		if err := r.scheduler.Delete(ctx, item.ConversationID, item, now); err != nil {
			return xerrors.Errorf("failed to create a schedule: %w", err)
		}
	}

	return nil
}

func (r *ReminderImpl) ListAll(ctx context.Context) (model.ReminderItems, error) {
	ctx, span := tracer.Start(ctx, "Reminder#ListAll")
	defer span.End()

	items, err := r.reminder.ListAll(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to list all reminder items: %w", err)
	}
	return items, nil
}

func (r *ReminderImpl) SyncSchedule(ctx context.Context, items model.ReminderItems) error {
	ctx, span := tracer.Start(ctx, "Reminder#SyncSchedule")
	defer span.End()

	now := time.Now()

	dl := logger.Default(ctx)
	dl.Info("start to sync schedule",
		zap.Any("items", items),
		zap.Int("count", len(items)),
	)

	items = items.FilterNextSchedule(now, syncInterval)

	dl.Info("FilterNextSchedule", zap.Int("count", len(items)))

	var conversationID model.ConversationID
	start := 0
	for i, item := range items {
		if conversationID == "" {
			conversationID = item.ConversationID
		} else if conversationID != item.ConversationID {
			if err := r.scheduler.Sync(ctx, conversationID, items[start:i], now); err != nil {
				return xerrors.Errorf("failed to sync schedule: %w", err)
			}
			conversationID = item.ConversationID
			start = i
		}
	}

	if len(items[start:]) > 0 {
		if err := r.scheduler.Sync(ctx, conversationID, items[start:], now); err != nil {
			return xerrors.Errorf("failed to sync schedule: %w", err)
		}
	}

	return nil
}
