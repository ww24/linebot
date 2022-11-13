package interactor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/internal/config"
)

const (
	triggerReminder = "リマインダー"
	prefixReminder  = "【リマインダー】"

	reminderDeletePrefix        = "Reminder#delete#"
	reminderDeleteConfirmPrefix = "Reminder#delete#confirm#"
)

type Reminder struct {
	conversation   service.Conversation
	reminder       service.Reminder
	message        repository.MessageProviderSet
	bot            service.Bot
	loc            *time.Location
	timeZoneOffset time.Duration
}

func NewReminder(
	conversation service.Conversation,
	reminder service.Reminder,
	message repository.MessageProviderSet,
	bot service.Bot,
	conf *config.Time,
) *Reminder {
	loc := conf.DefaultLocation()
	_, offset := time.Time{}.In(loc).Zone()

	return &Reminder{
		conversation:   conversation,
		reminder:       reminder,
		message:        message,
		bot:            bot,
		loc:            conf.DefaultLocation(),
		timeZoneOffset: time.Duration(offset) * time.Second,
	}
}

func (r *Reminder) Handle(ctx context.Context, e *model.Event) error {
	err := e.HandleTypeMessage(ctx, func(context.Context, *model.Event) error {
		if e.FilterText(triggerReminder) {
			return r.handleMenu(ctx, e)
		}

		return nil
	})
	if err != nil {
		return xerrors.Errorf("failed to handle type message: %w", err)
	}

	if err := e.HandleTypePostback(ctx, r.handlePostBack); err != nil {
		return xerrors.Errorf("failed to handle type postback: %w", err)
	}

	return nil
}

func (r *Reminder) handleMenu(ctx context.Context, e *model.Event) error {
	items, err := r.reminder.List(ctx, e.ConversationID())
	if err != nil {
		return xerrors.Errorf("failed to list reminder items: %w", err)
	}

	if len(items) == 0 {
		text := prefixReminder + "登録されていません。\n何をしますか？"
		msg := r.message.ReminderMenu(text, model.ReminderReplyTypeEmptyList, nil)
		if err := r.bot.ReplyMessage(ctx, e, msg); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
		return errResponseReturned
	}

	test := fmt.Sprintf(prefixReminder+"%d件登録されています。\n%s\n\n何をしますか？",
		len(items), items.Print(model.ListTypeOrdered))
	msg := r.message.ReminderMenu(test, model.ReminderReplyTypeAll, items)
	if err := r.bot.ReplyMessage(ctx, e, msg); err != nil {
		return xerrors.Errorf("failed to reply message: %w", err)
	}

	return errResponseReturned
}

func (r *Reminder) handlePostBack(ctx context.Context, e *model.Event) error {
	conversationID := e.ConversationID()

	switch e.Postback.Data {
	case "Reminder#add":
		status := &model.ConversationStatus{
			ConversationID: conversationID,
			Type:           model.ConversationStatusTypeReminderAdd,
		}
		if err := r.conversation.SetStatus(ctx, status); err != nil {
			return xerrors.Errorf("failed to set status: %w", err)
		}
		text := prefixReminder + "新規追加します。\n何をリマインドしますか？"
		msg := r.message.ReminderChoices(text,
			[]string{"買い物リスト"}, []model.ExecutorType{model.ExecutorTypeShoppingList})
		if err := r.bot.ReplyMessage(ctx, e, msg); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
		return errResponseReturned

	case "Reminder#add#shopping_list":
		text := prefixReminder + "買い物リストをリマインドします。\n何時にリマインドしますか？"
		msg := r.message.TimePicker(text, "Reminder#add#shopping_list#datetime")
		if err := r.bot.ReplyMessage(ctx, e, msg); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
		return errResponseReturned

	case "Reminder#add#shopping_list#datetime":
		t, err := time.Parse("15:04", e.Postback.Params.Time)
		if err != nil {
			return xerrors.Errorf("failed to parse time: %w", err)
		}
		t = t.In(r.loc).Add(-r.timeZoneOffset)
		text := prefixReminder + "毎日" + t.Format("15:04") + "に買い物リストをリマインドします。"
		if err := r.bot.ReplyMessage(ctx, e, r.message.Text(text)); err != nil {
			return xerrors.Errorf("failed to reply text message: %w", err)
		}
		item := model.NewReminderItem(
			conversationID,
			&model.DailyScheduler{
				Time: t,
			},
			&model.Executor{
				Type: model.ExecutorTypeShoppingList,
			},
		)
		if err := r.reminder.Add(ctx, item); err != nil {
			return xerrors.Errorf("failed to add reminder item: %w", err)
		}
		status := &model.ConversationStatus{
			ConversationID: conversationID,
			Type:           model.ConversationStatusTypeNeutral,
		}
		if err := r.conversation.SetStatus(ctx, status); err != nil {
			return xerrors.Errorf("failed to set status: %w", err)
		}
		return errResponseReturned
	}

	if err := r.handleDelete(ctx, e); err != nil {
		return err
	}

	return nil
}

func (r *Reminder) handleDelete(ctx context.Context, e *model.Event) error {
	switch {
	case strings.HasPrefix(e.Postback.Data, reminderDeleteConfirmPrefix):
		id := strings.TrimPrefix(e.Postback.Data, reminderDeleteConfirmPrefix)
		if err := r.reminder.Delete(ctx, e.ConversationID(), model.ReminderItemID(id)); err != nil {
			return xerrors.Errorf("failed to delete reminder item: %w", err)
		}
		msg := r.message.Text(prefixReminder + "リマインダーを削除しました。")
		if err := r.bot.ReplyMessage(ctx, e, msg); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
		return errResponseReturned

	case strings.HasPrefix(e.Postback.Data, reminderDeletePrefix):
		id := strings.TrimPrefix(e.Postback.Data, reminderDeletePrefix)
		msg := r.message.ReminderDeleteConfirmation("リマインダーを削除しますか？", reminderDeleteConfirmPrefix+id)
		if err := r.bot.ReplyMessage(ctx, e, msg); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
		return errResponseReturned
	}

	return nil
}

func (r *Reminder) HandleSchedule(ctx context.Context) error {
	items, err := r.reminder.ListAll(ctx)
	if err != nil {
		return xerrors.Errorf("failed to list reminder items: %w", err)
	}

	if err := r.reminder.SyncSchedule(ctx, items); err != nil {
		return xerrors.Errorf("failed to sync schedule: %w", err)
	}

	return nil
}
