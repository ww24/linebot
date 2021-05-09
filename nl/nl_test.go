package nl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()
	p, err := NewParser()
	require.NoError(t, err)

	tests := []struct {
		src  string
		want *Item
	}{
		{
			src: "りんごを削除",
			want: &Item{
				Name:   []string{"りんご"},
				Action: ActionTypeDelete,
			},
		},
		{
			src: "りんごを削除する。",
			want: &Item{
				Name:   []string{"りんご"},
				Action: ActionTypeDelete,
			},
		},
		{
			src: "りんごを消してください。",
			want: &Item{
				Name:   []string{"りんご"},
				Action: ActionTypeDelete,
			},
		},
		{
			src: "りんご削除",
			want: &Item{
				Name:   []string{"りんご"},
				Action: ActionTypeDelete,
			},
		},
		{
			src: "りんごを除去",
			want: &Item{
				Name:   []string{"りんご"},
				Action: ActionTypeDelete,
			},
		},
		{
			src: "1番を削除",
			want: &Item{
				Indexes: []int{1},
				Action:  ActionTypeDelete,
			},
		},
		{
			src: "１番目を削除。",
			want: &Item{
				Indexes: []int{1},
				Action:  ActionTypeDelete,
			},
		},
		{
			src: "１を削除",
			want: &Item{
				Indexes: []int{1},
				Action:  ActionTypeDelete,
			},
		},
		{
			src: "①を削除",
			want: &Item{
				Indexes: []int{1},
				Action:  ActionTypeDelete,
			},
		},
		{
			src: "1,2,3を削除.",
			want: &Item{
				Indexes: []int{1, 2, 3},
				Action:  ActionTypeDelete,
			},
		},
		{
			src: "１、２、３を消す。",
			want: &Item{
				Indexes: []int{1, 2, 3},
				Action:  ActionTypeDelete,
			},
		},
		{
			src: "1と2と3を除去",
			want: &Item{
				Indexes: []int{1, 2, 3},
				Action:  ActionTypeDelete,
			},
		},
		{
			src: "1と2と3を除去",
			want: &Item{
				Indexes: []int{1, 2, 3},
				Action:  ActionTypeDelete,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got := p.Parse(tt.src)
			assert.Equal(t, tt.want, got)
		})
	}
}
