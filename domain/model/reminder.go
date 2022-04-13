package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
)

func NewReminderItem(conversationID ConversationID, scheduler Scheduler, executor *Executor) *ReminderItem {
	return &ReminderItem{
		ID:             ReminderItemID(xid.New().String()),
		ConversationID: conversationID,
		Scheduler:      scheduler,
		Executor:       executor,
	}
}

type ReminderItemID string

type ReminderItem struct {
	ID             ReminderItemID
	ConversationID ConversationID
	Scheduler      Scheduler
	Executor       *Executor
}

type ReminderItemIDJSON struct {
	ConversationID string `json:"conversation_id"`
	ItemID         string `json:"item_id"`
}

func (r ReminderItem) IDJSON() *ReminderItemIDJSON {
	return &ReminderItemIDJSON{
		ConversationID: string(r.ConversationID),
		ItemID:         string(r.ID),
	}
}

type Executor struct {
	Type ExecutorType
}

type ExecutorType int

const (
	ExecutorTypeShoppingList ExecutorType = iota + 1
)

func (t ExecutorType) String() string {
	switch t {
	case ExecutorTypeShoppingList:
		return "shopping_list"
	default:
		return "unknown"
	}
}

func (t ExecutorType) UIText() string {
	switch t {
	case ExecutorTypeShoppingList:
		return "買い物リスト"
	default:
		return "unknown"
	}
}

type ReminderItems []*ReminderItem

func (l ReminderItems) Print(typ ListType) string {
	var b strings.Builder
	for i, item := range l {
		switch typ {
		case ListTypeOrdered:
			fmt.Fprint(&b, strconv.Itoa(i+1)+". ")

		case ListTypeDotted:
			fmt.Fprint(&b, "・")
		}
		fmt.Fprintf(&b, "%s: %s\n", item.Executor.Type, item.Scheduler)
	}
	return strings.TrimRight(b.String(), "\n")
}

func (l ReminderItems) FilterNextSchedule(t time.Time, d time.Duration) ReminderItems {
	threshold := t.Add(d)
	matched := make([]*ReminderItem, 0, len(l))
	for _, item := range l {
		next, err := item.Scheduler.Next(t)
		if err == nil && next.Before(threshold) {
			matched = append(matched, item)
		}
	}
	return matched
}
