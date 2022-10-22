package firestore

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/internal/code"
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
	ctx, span := r.tracer.Start(ctx, "Reminder#Add")
	defer span.End()

	entity := NewReminderItem(item)
	entity.CreatedAt = r.now().Unix()

	reminder := r.reminder(item.ConversationID)
	if _, err := reminder.Doc(entity.ID).Set(ctx, entity); err != nil {
		return xerrors.Errorf("failed to add reminder: %w", err)
	}

	return nil
}

func (r *Reminder) List(ctx context.Context, conversationID model.ConversationID) ([]*model.ReminderItem, error) {
	ctx, span := r.tracer.Start(ctx, "Reminder#List")
	defer span.End()

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

		m, err := item.Model(conversationID, doc.Ref.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, m)
	}

	return items, nil
}

func (r *Reminder) Get(ctx context.Context, conversationID model.ConversationID, itemID model.ReminderItemID) (*model.ReminderItem, error) {
	ctx, span := r.tracer.Start(ctx, "Reminder#Get")
	defer span.End()

	doc, err := r.reminder(conversationID).Doc(string(itemID)).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			err = code.With(err, code.NotFound)
		}
		return nil, xerrors.Errorf("failed to get reminder: %w", err)
	}

	var item ReminderItem
	if err := doc.DataTo(&item); err != nil {
		return nil, xerrors.Errorf("failed to convert data to item: %w", err)
	}

	m, err := item.Model(conversationID, doc.Ref.ID)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *Reminder) Delete(ctx context.Context, conversationID model.ConversationID, id model.ReminderItemID) error {
	ctx, span := r.tracer.Start(ctx, "Reminder#Delete")
	defer span.End()

	reminder := r.reminder(conversationID)
	if _, err := reminder.Doc(string(id)).Delete(ctx); err != nil {
		return xerrors.Errorf("failed to delete reminder: %w", err)
	}

	return nil
}

func (r *Reminder) ListAll(ctx context.Context) ([]*model.ReminderItem, error) {
	ctx, span := r.tracer.Start(ctx, "Reminder#ListAll")
	defer span.End()

	// TODO: return wrapped iterator for performance

	items := make([]*model.ReminderItem, 0)
	iter := r.conversations().Documents(ctx)
	for {
		ss, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, xerrors.Errorf("failed to iterate conversations: %w", err)
		}

		conversationID := model.ConversationID(ss.Ref.ID)
		docs, err := ss.Ref.Collection("reminders").OrderBy("created_at", firestore.Asc).Documents(ctx).GetAll()
		if err != nil {
			return nil, xerrors.Errorf("failed to find reminder: %w", err)
		}

		for _, doc := range docs {
			var item ReminderItem
			if err := doc.DataTo(&item); err != nil {
				return nil, xerrors.Errorf("failed to convert data to item: %w", err)
			}

			m, err := item.Model(conversationID, doc.Ref.ID)
			if err != nil {
				return nil, err
			}
			items = append(items, m)
		}
	}

	return items, nil
}

type ReminderItem struct {
	ConversationID model.ConversationID `firestore:"-"`
	ID             string               `firestore:"-"`
	Scheduler      string               `firestore:"scheduler"`
	Executor       *Executor            `firestore:"executor"`
	CreatedAt      int64                `firestore:"created_at"` // UNIX time
}

type Executor struct {
	Type model.ExecutorType `firestore:"type"`
}

func NewReminderItem(src *model.ReminderItem) *ReminderItem {
	return &ReminderItem{
		ConversationID: src.ConversationID,
		ID:             string(src.ID),
		Scheduler:      src.Scheduler.String(),
		Executor:       NewExecutor(src.Executor),
	}
}

func (r *ReminderItem) Model(conversationID model.ConversationID, id string) (*model.ReminderItem, error) {
	sch, err := model.ParseScheduler(r.Scheduler)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse scheduler: %w", err)
	}

	return &model.ReminderItem{
		ConversationID: conversationID,
		ID:             model.ReminderItemID(id),
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
