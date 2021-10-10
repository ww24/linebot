package model

import (
	"reflect"
	"testing"
)

func TestItem_UniqueIndexes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		item *Item
		want []int
	}{
		{
			name: "nil indexes",
			item: &Item{
				Indexes: nil,
			},
			want: []int{},
		},
		{
			name: "empty indexes",
			item: &Item{
				Indexes: []int{},
			},
			want: []int{},
		},
		{
			name: "sequential indexes",
			item: &Item{
				Indexes: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "duplicate indexes",
			item: &Item{
				Indexes: []int{1, 1, 2, 3, 3, 5},
			},
			want: []int{1, 2, 3, 5},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.item.UniqueIndexes()
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got: %+v, want: %+v", got, tt.want)
			}
		})
	}
}
