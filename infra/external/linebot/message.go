package linebot

import (
	"errors"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
)

var (
	errInvalidMessageType = errors.New("invalid given message type")
)

// MessageProviderSet implements repository.MessageProviderSet.
type MessageProviderSet struct{}

func NewMessageProviderSet() *MessageProviderSet {
	return &MessageProviderSet{}
}

func (s *MessageProviderSet) Text(text string) repository.MessageProvider {
	return &TextMessage{text: text}
}

func (s *MessageProviderSet) ShoppingDeleteConfirmation(text string) repository.MessageProvider {
	return &ShoppingDeleteConfirmation{text: text}
}

func (s *MessageProviderSet) ShoppingMenu(text string, rt model.ShoppingReplyType) repository.MessageProvider {
	return &ShoppingMenu{
		text:      text,
		replyType: rt,
	}
}

func (s *MessageProviderSet) ReminderMenu(text string, rt model.ReminderReplyType, items []*model.ReminderItem) repository.MessageProvider {
	reminderMenu := &ReminderMenu{
		text:      text,
		replyType: rt,
	}
	t := time.Now()
	if len(items) == 0 {
		return reminderMenu
	}

	if data, err := makeReminderListMessage(items, t); err == nil {
		if flexContainer, err := linebot.UnmarshalFlexMessageJSON(data); err == nil {
			reminderMenu.flex = &flexContainer
		}
	}
	return reminderMenu
}

func (s *MessageProviderSet) ReminderChoices(text string, labels []string, types []model.ExecutorType) repository.MessageProvider {
	return &ReminderChoices{
		text:   text,
		labels: labels,
		types:  types,
	}
}

func (s *MessageProviderSet) TimePicker(text, data string) repository.MessageProvider {
	return &TimePicker{
		text: text,
		data: data,
	}
}

func (s *MessageProviderSet) ReminderDeleteConfirmation(text, data string) repository.MessageProvider {
	return &ReminderDeleteConfirmation{
		text: text,
		data: data,
	}
}

func asMessage(msg linebot.SendingMessage, v interface{}) error {
	m, ok := v.(*linebot.SendingMessage)
	if !ok {
		return xerrors.Errorf("invalid message: %w", errInvalidMessageType)
	}

	*m = msg

	return nil
}

type TextMessage struct {
	text string
}

func (p *TextMessage) AsMessage(v interface{}) error {
	msg := linebot.NewTextMessage(p.text)
	return asMessage(msg, v)
}

// Message implements repository.MessageProvider.
type ShoppingDeleteConfirmation struct {
	text string
}

// AsMessage assigns the message to the given *linebot.SendingMessage.
func (p *ShoppingDeleteConfirmation) AsMessage(v interface{}) error {
	var msg linebot.SendingMessage
	msg = linebot.NewTextMessage(p.text)
	msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
		Items: []*linebot.QuickReplyButton{
			{Action: linebot.NewPostbackAction("YES", "Shopping#deleteConfirm", "", "YES")},
			{Action: linebot.NewPostbackAction("NO", "Shopping#deleteCancel", "", "NO")},
		},
	})

	return asMessage(msg, v)
}

// ShoppingMenu implements repository.MessageProvider.
type ShoppingMenu struct {
	text      string
	replyType model.ShoppingReplyType
}

func (p *ShoppingMenu) AsMessage(v interface{}) error {
	var msg linebot.SendingMessage
	msg = linebot.NewTextMessage(p.text)

	//nolint: exhaustive
	switch p.replyType {
	case model.ShoppingReplyTypeEmptyList:
		msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
			Items: []*linebot.QuickReplyButton{
				{Action: linebot.NewPostbackAction("追加", "Shopping#add", "", "追加")},
			},
		})
	case model.ShoppingReplyTypeWithoutView:
		msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
			Items: []*linebot.QuickReplyButton{
				{Action: linebot.NewPostbackAction("削除", "Shopping#delete", "", "削除")},
				{Action: linebot.NewPostbackAction("追加", "Shopping#add", "", "追加")},
			},
		})
	default:
		msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
			Items: []*linebot.QuickReplyButton{
				{Action: linebot.NewPostbackAction("削除", "Shopping#delete", "", "削除")},
				{Action: linebot.NewPostbackAction("追加", "Shopping#add", "", "追加")},
				{Action: linebot.NewPostbackAction("表示", "Shopping#view", "", "表示")},
			},
		})
	}

	return asMessage(msg, v)
}

// ShoppingMenu implements repository.MessageProvider.
type ReminderMenu struct {
	text      string
	flex      *linebot.FlexContainer
	replyType model.ReminderReplyType
}

func (r *ReminderMenu) AsMessage(v interface{}) error {
	var msg linebot.SendingMessage
	if r.flex != nil {
		msg = linebot.NewFlexMessage(r.text, *r.flex)
	} else {
		msg = linebot.NewTextMessage(r.text)
	}

	//nolint: exhaustive
	switch r.replyType {
	default:
		msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
			Items: []*linebot.QuickReplyButton{
				{Action: linebot.NewPostbackAction("追加", "Reminder#add", "", "追加")},
			},
		})
	}

	return asMessage(msg, v)
}

type ReminderChoices struct {
	text   string
	labels []string
	types  []model.ExecutorType
}

func (r *ReminderChoices) AsMessage(v interface{}) error {
	items := make([]*linebot.QuickReplyButton, 0, len(r.labels))
	for i := range r.labels {
		label := r.labels[i]
		items = append(items, &linebot.QuickReplyButton{
			Action: linebot.NewPostbackAction(label, "Reminder#add#"+r.types[i].String(), "", label),
		})
	}

	var msg linebot.SendingMessage
	msg = linebot.NewTextMessage(r.text)
	msg = msg.WithQuickReplies(&linebot.QuickReplyItems{Items: items})

	return asMessage(msg, v)
}

type TimePicker struct {
	text string
	data string
}

func (p *TimePicker) AsMessage(v interface{}) error {
	var msg linebot.SendingMessage
	msg = linebot.NewTextMessage(p.text)
	msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
		Items: []*linebot.QuickReplyButton{
			{Action: linebot.NewDatetimePickerAction("時刻設定", p.data, "time", "", "", "")},
		},
	})

	return asMessage(msg, v)
}

type ReminderDeleteConfirmation struct {
	text string
	data string
}

func (c *ReminderDeleteConfirmation) AsMessage(v interface{}) error {
	var msg linebot.SendingMessage
	msg = linebot.NewTextMessage(c.text)
	msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
		Items: []*linebot.QuickReplyButton{
			{Action: linebot.NewPostbackAction("YES", c.data, "", "YES")},
			{Action: linebot.NewPostbackAction("NO", "Reminder#cancel", "", "NO")},
		},
	})

	return asMessage(msg, v)
}
