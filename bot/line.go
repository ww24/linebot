package bot

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/ww24/linebot/domain/model"
)

func ConversationID(source *linebot.EventSource) model.ConversationID {
	switch source.Type {
	case linebot.EventSourceTypeGroup:
		return model.NewConversationID("LG", source.GroupID)
	case linebot.EventSourceTypeRoom:
		return model.NewConversationID("LR", source.RoomID)
	case linebot.EventSourceTypeUser:
		return model.NewConversationID("LU", source.UserID)
	default:
		return model.NewConversationID("LX", source.UserID)
	}
}
