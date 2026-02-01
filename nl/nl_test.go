package nl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ww24/linebot/domain/model"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()
	p, err := NewParser()
	require.NoError(t, err)

	tests := []struct {
		src  string
		want *model.Item
	}{
		{
			src: "りんごを削除",
			want: &model.Item{
				Name:   []string{"りんご"},
				Action: model.ActionTypeDelete,
			},
		},
		{
			src: "りんごを削除する。",
			want: &model.Item{
				Name:   []string{"りんご"},
				Action: model.ActionTypeDelete,
			},
		},
		{
			src: "りんごを消してください。",
			want: &model.Item{
				Name:   []string{"りんご"},
				Action: model.ActionTypeDelete,
			},
		},
		{
			src: "りんごを除去",
			want: &model.Item{
				Name:   []string{"りんご"},
				Action: model.ActionTypeDelete,
			},
		},
		{
			src: "1番を削除",
			want: &model.Item{
				Indexes: []int{1},
				Action:  model.ActionTypeDelete,
			},
		},
		{
			src: "１番目を削除。",
			want: &model.Item{
				Indexes: []int{1},
				Action:  model.ActionTypeDelete,
			},
		},
		{
			src: "１を削除",
			want: &model.Item{
				Indexes: []int{1},
				Action:  model.ActionTypeDelete,
			},
		},
		{
			src: "①を削除",
			want: &model.Item{
				Indexes: []int{1},
				Action:  model.ActionTypeDelete,
			},
		},
		{
			src: "1,2,3を削除.",
			want: &model.Item{
				Indexes: []int{1, 2, 3},
				Action:  model.ActionTypeDelete,
			},
		},
		{
			src: "１、２、３を消す。",
			want: &model.Item{
				Indexes: []int{1, 2, 3},
				Action:  model.ActionTypeDelete,
			},
		},
		{
			src: "1と2と3を除去",
			want: &model.Item{
				Indexes: []int{1, 2, 3},
				Action:  model.ActionTypeDelete,
			},
		},
		{
			src: "1と2と3を除去",
			want: &model.Item{
				Indexes: []int{1, 2, 3},
				Action:  model.ActionTypeDelete,
			},
		},
		{
			src: "11を削除",
			want: &model.Item{
				Indexes: []int{11},
				Action:  model.ActionTypeDelete,
			},
		},
		{
			src: "11,12を削除",
			want: &model.Item{
				Indexes: []int{11, 12},
				Action:  model.ActionTypeDelete,
			},
		},
		{
			src: "11と12を削除",
			want: &model.Item{
				Indexes: []int{11, 12},
				Action:  model.ActionTypeDelete,
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got := p.Parse(tt.src)
			assert.Equal(t, tt.want, got)
		})
	}
}
