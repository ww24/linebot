package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ww24/linebot/domain/model"
	"golang.org/x/xerrors"
)

type Reminder struct {
	*Conversation
}

func NewReminder(c *Conversation) *Reminder {
	return &Reminder{Conversation: c}
}

func (r *Reminder) reminder(conversationID model.ConversationID) *firestore.CollectionRef {
	return r.conversation(conversationID).Collection("reminders")
}

func (r *Reminder) Add(ctx context.Context, item *model.ReminderItem) error {
	entity := NewReminderItem(item)
	entity.CreatedAt = time.Now().Unix()

	reminder := r.reminder(item.ConversationID)
	if _, err := reminder.NewDoc().Set(ctx, entity); err != nil {
		return xerrors.Errorf("failed to add reminder: %w", err)
	}

	return nil
}

func (r *Reminder) Find(ctx context.Context, conversationID model.ConversationID) ([]*model.ReminderItem, error) {
	reminder := r.reminder(conversationID)
	iter := reminder.OrderBy("created_at", firestore.Asc).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, xerrors.Errorf("failed to find reminder: %w", err)
	}

	items := make([]*model.ReminderItem, 0, len(docs))
	for _, doc := range docs {
		var item ReminderItem
		if err := doc.DataTo(&item); err != nil {
			return nil, xerrors.Errorf("failed to convert data to item: %w", err)
		}

		item.ID = doc.Ref.ID
		m, err := item.Model()
		if err != nil {
			return nil, err
		}
		items = append(items, m)
	}

	return items, nil
}

func (r *Reminder) Delete(ctx context.Context, conversationID model.ConversationID, id string) error {
	reminder := r.reminder(conversationID)
	if _, err := reminder.Doc(id).Delete(ctx); err != nil {
		return xerrors.Errorf("failed to delete reminder: %w", err)
	}

	return nil
}

type ReminderItem struct {
	ID             string    `firestore:"-"`
	Name           string    `firestore:"name"`
	ConversationID string    `firestore:"conversation_id"`
	Scheduler      string    `firestore:"scheduler"`
	Executor       *Executor `firestore:"executor"`
	CreatedAt      int64     `firestore:"created_at"` // UNIX time
}

type Executor struct {
	Type model.ExecutorType `firestore:"type"`
}

func NewReminderItem(src *model.ReminderItem) *ReminderItem {
	return &ReminderItem{
		ID:             src.ID,
		Name:           src.Name,
		ConversationID: string(src.ConversationID),
		Scheduler:      src.Scheduler.String(),
		Executor:       NewExecutor(src.Executor),
	}
}

func (r *ReminderItem) Model() (*model.ReminderItem, error) {
	sch, err := model.ParseScheduler(r.Scheduler)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse scheduler: %w", err)
	}
	return &model.ReminderItem{
		ID:             r.ID,
		Name:           r.Name,
		ConversationID: model.ConversationID(r.ConversationID),
		Scheduler:      sch,
		Executor:       r.Executor.Model(),
	}, nil
}

func NewExecutor(src *model.Executor) *Executor {
	return &Executor{
		Type: src.Type,
	}
}

func (e *Executor) Model() *model.Executor {
	return &model.Executor{
		Type: e.Type,
	}
}
