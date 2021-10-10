package model

import (
	"errors"
	"strings"

	"golang.org/x/xerrors"
)

type ConversationStatusType int

const (
	ConversationStatusTypeNeutral ConversationStatusType = iota
	ConversationStatusTypeShopping
	ConversationStatusTypeShoppingAdd
)

const (
	conversationIDSep        = "_"
	conversationSeparateSize = 2
)

var errConversationStatusValidationFailed = errors.New("conversation status validation failed")

type ConversationID string

func NewConversationID(prefix, sourceID string) ConversationID {
	return ConversationID(prefix + conversationIDSep + sourceID)
}

func (c ConversationID) SourceID() string {
	s := strings.SplitN(string(c), conversationIDSep, conversationSeparateSize)
	if len(s) < conversationSeparateSize {
		return ""
	}
	return s[1]
}

func (c ConversationID) String() string {
	return string(c)
}

type ConversationStatus struct {
	ConversationID ConversationID
	Type           ConversationStatusType
}

func (m *ConversationStatus) Validate() error {
	if m.ConversationID == "" {
		return xerrors.Errorf("invalid empty conversation id: %w", errConversationStatusValidationFailed)
	}
	return nil
}
