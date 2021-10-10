package model

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

var (
	errShoppingItemValidationFailed = errors.New("shopping item validation failed")
)

type ShoppingItem struct {
	ID             string
	Name           string
	Quantity       int
	ConversationID ConversationID
	CreatedAt      int64
	Order          int
}

func (m *ShoppingItem) Validate() error {
	if m.ConversationID == "" {
		return xerrors.Errorf("invalid empty conversation id: %w", errShoppingItemValidationFailed)
	}
	if m.Name == "" {
		return xerrors.Errorf("invalid empty name: %w", errShoppingItemValidationFailed)
	}
	if m.CreatedAt == 0 {
		return xerrors.Errorf("invalid created at: %w", errShoppingItemValidationFailed)
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

		case ListTypeDotted:
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
