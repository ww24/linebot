package code

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestInternalError_From(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want Code
	}{
		{
			name: "nil",
			err:  nil,
			want: OK,
		},
		{
			name: "internalError",
			err: &internalError{
				source: errors.New("test"),
				code:   NotFound,
			},
			want: NotFound,
		},
		{
			name: "wrapped internalError",
			err: fmt.Errorf("wrapped: %w", &internalError{
				source: errors.New("test"),
				code:   NotFound,
			}),
			want: NotFound,
		},
		{
			name: "internalError wrapped by xerrors",
			err: xerrors.Errorf("wrapped: %w", &internalError{
				source: errors.New("test"),
				code:   NotFound,
			}),
			want: NotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := From(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
