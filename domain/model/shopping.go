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

type ListType int

const (
	ListTypeDotted ListType = iota
	ListTypeOrdered
)

func (l ShoppingItems) Print(typ ListType) string {
	var b strings.Builder
	for i, item := range l {
		switch typ {
		case ListTypeOrdered:
			fmt.Fprintf(&b, "%d. %s\n", i+1, item.Name)

		default:
			fmt.Fprintf(&b, "・%s\n", item.Name)
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func (l ShoppingItems) FilterByNames(names []string) ShoppingItems {
	res := make([]*ShoppingItem, 0)
	for _, item := range l {
		for _, name := range names {
			if strings.Contains(item.Name, name) {
				res = append(res, item)
				break
			}
		}
	}
	return res
}
