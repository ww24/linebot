package model

import (
	"errors"
	"fmt"
	"strings"
)

type ShoppingItem struct {
	ID             string
	Name           string
	Quantity       int
	ConversationID string
	CreatedAt      int64
}

func (m *ShoppingItem) Validate() error {
	if m.ConversationID == "" {
		return errors.New("invalid empty conversation id")
	}
	if m.Name == "" {
		return errors.New("invalid empty name")
	}
	if m.CreatedAt == 0 {
		return errors.New("invalid created at")
	}
	return nil
}

type ShoppingItems []*ShoppingItem

func (l ShoppingItems) Print() string {
	var b strings.Builder
	for i, item := range l {
		fmt.Fprintf(&b, "%d. %s\n", i+1, item.Name)
	}
	return strings.TrimRight(b.String(), "\n")
}
