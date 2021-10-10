package linebot

import (
	"errors"

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
