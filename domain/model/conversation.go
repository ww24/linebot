package model

import "errors"

type ConversationStatusType int

const (
	ConversationStatusTypeNeutral ConversationStatusType = iota
	ConversationStatusTypeShopping
	ConversationStatusTypeShoppingAdd
)

type ConversationStatus struct {
	ConversationID string
	Type           ConversationStatusType
}

func (m *ConversationStatus) Validate() error {
	if m.ConversationID == "" {
		return errors.New("invalid empty conversation id")
	}
	return nil
}
