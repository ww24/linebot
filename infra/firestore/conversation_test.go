package firestore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/internal/code"
)

func TestConversation_SetStatus(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := []struct {
		name    string
		status  *model.ConversationStatus
		wantErr error
	}{
		{
			name: "set conversation type neutral",
			status: &model.ConversationStatus{
				ConversationID: "conv_set_neutral",
				Type:           model.ConversationStatusTypeNeutral,
			},
			wantErr: nil,
		},
		{
			name: "set conversation type shopping",
			status: &model.ConversationStatus{
				ConversationID: "conv_set_shopping",
				Type:           model.ConversationStatusTypeShopping,
			},
			wantErr: nil,
		},
		{
			name: "try to set invalid conversation type",
			status: &model.ConversationStatus{
				ConversationID: "conv_set_invalid",
				Type:           model.ConversationStatusType(-1),
			},
			wantErr: model.ErrConversationStatusValidationFailed,
		},
	}
	conv := NewConversation(testCli)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := conv.SetStatus(ctx, tt.status)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr != nil {
				return
			}

			ss, err := testCli.cli.Collection("conversations").Doc(tt.status.ConversationID.String()).Get(ctx)
			require.NoError(t, err)
			var status ConversationStatus
			require.NoError(t, ss.DataTo(&status))
			assert.Equal(t, tt.status, status.Model(tt.status.ConversationID))
		})
	}
}

func TestConversation_GetStatus(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	testdata := []*ConversationStatus{
		{
			ConversationID: "get_conv_status",
			Status:         int(model.ConversationStatusTypeShopping),
		},
	}
	for _, data := range testdata {
		_, err := testCli.cli.Collection("conversations").Doc(data.ConversationID.String()).Set(ctx, data)
		require.NoError(t, err)
	}
	tests := []struct {
		name     string
		id       model.ConversationID
		want     *model.ConversationStatus
		wantCode code.Code
	}{
		{
			name: "get conversation status",
			id:   "get_conv_status",
			want: &model.ConversationStatus{
				ConversationID: "get_conv_status",
				Type:           model.ConversationStatusTypeShopping,
			},
			wantCode: code.OK,
		},
		{
			name:     "conversation status not found",
			id:       "conv_status_not_found",
			want:     nil,
			wantCode: code.NotFound,
		},
	}
	conv := NewConversation(testCli)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := conv.GetStatus(ctx, tt.id)
			assert.Equal(t, tt.wantCode, code.From(err), "Error: %#v", err)
			assert.Equal(t, tt.want, got)
		})
	}
}
