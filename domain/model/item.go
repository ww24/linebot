package model

type ActionType int

const (
	ActionTypeUnknown ActionType = iota
	ActionTypeDelete
)

type Item struct {
	Indexes []int
	Name    []string
	Action  ActionType
}
