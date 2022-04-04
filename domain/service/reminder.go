package service

import (
	"context"
	"time"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
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

func NewReminder(reminder repository.Reminder, scheduler repository.ScheduleSynchronizer) *ReminderImpl {
	return &ReminderImpl{
		reminder:  reminder,
		scheduler: scheduler,
	}
}

func (r *ReminderImpl) Add(ctx context.Context, reminder *model.ReminderItem) error {
	if err := r.reminder.Add(ctx, reminder); err != nil {
		return xerrors.Errorf("failed to add a reminder item: %w", err)
	}
	return nil
}

func (r *ReminderImpl) List(ctx context.Context, conversationID model.ConversationID) (model.ReminderItems, error) {
	items, err := r.reminder.List(ctx, conversationID)
	if err != nil {
		return nil, xerrors.Errorf("failed to list reminder items: %w", err)
	}
	return items, nil
}

func (r *ReminderImpl) Get(ctx context.Context, conversationID model.ConversationID, itemID model.ReminderItemID) (*model.ReminderItem, error) {
	item, err := r.reminder.Get(ctx, conversationID, itemID)
	if err != nil {
		return nil, xerrors.Errorf("failed to get a reminder item: %w", err)
	}
	return item, nil
}

func (r *ReminderImpl) Delete(ctx context.Context, conversationID model.ConversationID, reminderItemID model.ReminderItemID) error {
	if err := r.reminder.Delete(ctx, conversationID, reminderItemID); err != nil {
		return xerrors.Errorf("failed to delete a reminder item: %w", err)
	}
	return nil
}

func (r *ReminderImpl) ListAll(ctx context.Context) (model.ReminderItems, error) {
	items, err := r.reminder.ListAll(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to list all reminder items: %w", err)
	}
	return items, nil
}

func (r *ReminderImpl) SyncSchedule(ctx context.Context, items model.ReminderItems) error {
	now := time.Now()
	items = items.FilterNextSchedule(now, syncInterval)

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
