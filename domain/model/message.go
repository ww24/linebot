package model

type ShoppingReplyType int

const (
	ShoppingReplyTypeAll ShoppingReplyType = iota
	ShoppingReplyTypeEmptyList
	ShoppingReplyTypeWithoutView
)

type ReminderReplyType int

const (
	ReminderReplyTypeAll ReminderReplyType = iota
	ReminderReplyTypeEmptyList
)
