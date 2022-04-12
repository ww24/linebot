package model

import (
	"context"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Event struct {
	*linebot.Event
	Status *ConversationStatus
}

// ConversationID returns conversation ID.
func (e *Event) ConversationID() ConversationID {
	switch e.Source.Type {
	case linebot.EventSourceTypeGroup:
		return NewConversationID("LG", e.Source.GroupID)
	case linebot.EventSourceTypeRoom:
		return NewConversationID("LR", e.Source.RoomID)
	case linebot.EventSourceTypeUser:
		return NewConversationID("LU", e.Source.UserID)
	default:
		return NewConversationID("LX", e.Source.UserID)
	}
}

func (e *Event) SetStatus(st ConversationStatusType) {
	e.Status = &ConversationStatus{
		ConversationID: e.ConversationID(),
		Type:           st,
	}
}

func (e *Event) HandleTypeMessage(ctx context.Context, f func(context.Context, *Event) error) error {
	if e.Type == linebot.EventTypeMessage {
		return f(ctx, e)
	}
	return nil
}

func (e *Event) HandleTypePostback(ctx context.Context, f func(context.Context, *Event) error) error {
	if e.Type == linebot.EventTypePostback {
		return f(ctx, e)
	}
	return nil
}

// FilterText returns true if Event.Message contains target text.
func (e *Event) FilterText(target string) bool {
	text, ok := e.Message.(*linebot.TextMessage)
	if ok && strings.Contains(text.Text, target) {
		return true
	}

	return false
}

// ReadTextLines reads text lines from Event.Message and trim spaces per line.
func (e *Event) ReadTextLines() []string {
	text, ok := e.Message.(*linebot.TextMessage)
	if !ok {
		return nil
	}

	lines := strings.Split(text.Text, "\n")
	ret := make([]string, 0, len(lines))
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			ret = append(ret, line)
		}
	}

	return ret
}
