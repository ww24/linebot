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

func (i *Item) UniqueIndexes() []int {
	indexes := make([]int, 0)
	set := make(map[int]struct{}, len(i.Indexes))
	for _, index := range i.Indexes {
		if _, ok := set[index]; !ok {
			indexes = append(indexes, index)
			set[index] = struct{}{}
		}
	}
	return indexes
}
