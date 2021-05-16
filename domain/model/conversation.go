package model

import (
	"errors"
	"strings"
)

type ConversationStatusType int

const (
	ConversationStatusTypeNeutral ConversationStatusType = iota
	ConversationStatusTypeShopping
	ConversationStatusTypeShoppingAdd
)

const conversationIDSep = "_"

type ConversationID string

func NewConversationID(prefix, sourceID string) ConversationID {
	return ConversationID(prefix + conversationIDSep + sourceID)
}

func (c ConversationID) SourceID() string {
	s := strings.SplitN(string(c), conversationIDSep, 2)
	if len(s) < 2 {
		return ""
	}
	return s[1]
}

type ConversationStatus struct {
	ConversationID ConversationID
	Type           ConversationStatusType
}

func (m *ConversationStatus) Validate() error {
	if m.ConversationID == "" {
		return errors.New("invalid empty conversation id")
	}
	return nil
}
