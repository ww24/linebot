package linebot

import (
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
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

func (s *MessageProviderSet) Image(originalURL, previewURL string) repository.MessageProvider {
	return &Image{
		originalURL: originalURL,
		previewURL:  previewURL,
	}
}

type TextMessage struct {
	text string
}

func (p *TextMessage) ToMessage() linebot.SendingMessage {
	return linebot.NewTextMessage(p.text)
}

// Message implements repository.MessageProvider.
type ShoppingDeleteConfirmation struct {
	text string
}

// AsMessage assigns the message to the given *linebot.SendingMessage.
func (p *ShoppingDeleteConfirmation) ToMessage() linebot.SendingMessage {
	var msg linebot.SendingMessage
	msg = linebot.NewTextMessage(p.text)
	msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
		Items: []*linebot.QuickReplyButton{
			{Action: linebot.NewPostbackAction("YES", "Shopping#deleteConfirm", "", "YES", "", "")},
			{Action: linebot.NewPostbackAction("NO", "Shopping#deleteCancel", "", "NO", "", "")},
		},
	})
	return msg
}

// ShoppingMenu implements repository.MessageProvider.
type ShoppingMenu struct {
	text      string
	replyType model.ShoppingReplyType
}

func (p *ShoppingMenu) ToMessage() linebot.SendingMessage {
	var msg linebot.SendingMessage
	msg = linebot.NewTextMessage(p.text)

	switch p.replyType {
	case model.ShoppingReplyTypeEmptyList:
		msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
			Items: []*linebot.QuickReplyButton{
				{Action: linebot.NewPostbackAction("追加", "Shopping#add", "", "追加", "", "")},
			},
		})
	case model.ShoppingReplyTypeWithoutView:
		msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
			Items: []*linebot.QuickReplyButton{
				{Action: linebot.NewPostbackAction("削除", "Shopping#delete", "", "削除", "", "")},
				{Action: linebot.NewPostbackAction("追加", "Shopping#add", "", "追加", "", "")},
			},
		})
	default:
		msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
			Items: []*linebot.QuickReplyButton{
				{Action: linebot.NewPostbackAction("削除", "Shopping#delete", "", "削除", "", "")},
				{Action: linebot.NewPostbackAction("追加", "Shopping#add", "", "追加", "", "")},
				{Action: linebot.NewPostbackAction("表示", "Shopping#view", "", "表示", "", "")},
			},
		})
	}

	return msg
}

// ShoppingMenu implements repository.MessageProvider.
type ReminderMenu struct {
	text      string
	flex      *linebot.FlexContainer
	replyType model.ReminderReplyType
}

func (r *ReminderMenu) ToMessage() linebot.SendingMessage {
	var msg linebot.SendingMessage
	if r.flex != nil {
		msg = linebot.NewFlexMessage(r.text, *r.flex)
	} else {
		msg = linebot.NewTextMessage(r.text)
	}

	return msg.WithQuickReplies(&linebot.QuickReplyItems{
		Items: []*linebot.QuickReplyButton{
			{Action: linebot.NewPostbackAction("追加", "Reminder#add", "", "追加", "", "")},
		},
	})
}

type ReminderChoices struct {
	text   string
	labels []string
	types  []model.ExecutorType
}

func (r *ReminderChoices) ToMessage() linebot.SendingMessage {
	items := make([]*linebot.QuickReplyButton, 0, len(r.labels))
	for i := range r.labels {
		label := r.labels[i]
		items = append(items, &linebot.QuickReplyButton{
			Action: linebot.NewPostbackAction(label, "Reminder#add#"+r.types[i].String(), "", label, "", ""),
		})
	}

	var msg linebot.SendingMessage
	msg = linebot.NewTextMessage(r.text)
	msg = msg.WithQuickReplies(&linebot.QuickReplyItems{Items: items})

	return msg
}

type TimePicker struct {
	text string
	data string
}

func (p *TimePicker) ToMessage() linebot.SendingMessage {
	var msg linebot.SendingMessage
	msg = linebot.NewTextMessage(p.text)
	msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
		Items: []*linebot.QuickReplyButton{
			{Action: linebot.NewDatetimePickerAction("時刻設定", p.data, "time", "", "", "")},
		},
	})

	return msg
}

type ReminderDeleteConfirmation struct {
	text string
	data string
}

func (c *ReminderDeleteConfirmation) ToMessage() linebot.SendingMessage {
	var msg linebot.SendingMessage
	msg = linebot.NewTextMessage(c.text)
	msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
		Items: []*linebot.QuickReplyButton{
			{Action: linebot.NewPostbackAction("YES", c.data, "", "YES", "", "")},
			{Action: linebot.NewPostbackAction("NO", "Reminder#cancel", "", "NO", "", "")},
		},
	})

	return msg
}

type Image struct {
	originalURL string
	previewURL  string
}

func (i *Image) ToMessage() linebot.SendingMessage {
	return linebot.NewImageMessage(i.originalURL, i.previewURL)
}
