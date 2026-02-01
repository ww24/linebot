package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversationStatus_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status *ConversationStatus
		want   error
	}{
		{
			name: "success",
			status: &ConversationStatus{
				ConversationID: "conv_id",
				Type:           ConversationStatusTypeNeutral,
			},
			want: nil,
		},
		{
			name: "empty id",
			status: &ConversationStatus{
				ConversationID: "",
				Type:           ConversationStatusTypeNeutral,
			},
			want: ErrConversationStatusValidationFailed,
		},
		{
			name: "invalid conversation status",
			status: &ConversationStatus{
				ConversationID: "invalid",
				Type:           ConversationStatusType(4),
			},
			want: ErrConversationStatusValidationFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.status.Validate()
			assert.ErrorIs(t, err, tt.want)
		})
	}
}
