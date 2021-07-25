package model

type ReminderItem struct {
	ID             string
	Name           string
	ConversationID ConversationID
	Scheduler      Scheduler
	Executor       *Executor
}

type Executor struct {
	Type ExecutorType
}

type ExecutorType int

const (
	ExecutorTypeShoppingList = iota + 1
)
