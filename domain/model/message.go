package model

type ShoppingReplyType int

const (
	ShoppingReplyTypeAll ShoppingReplyType = iota
	ShoppingReplyTypeEmptyList
	ShoppingReplyTypeWithoutView
)
