package bot

import "github.com/line/line-bot-sdk-go/v7/linebot"

func ConversationID(source *linebot.EventSource) string {
	switch source.Type {
	case linebot.EventSourceTypeGroup:
		return "LG_" + source.GroupID
	case linebot.EventSourceTypeRoom:
		return "LR_" + source.RoomID
	case linebot.EventSourceTypeUser:
		return "LU_" + source.UserID
	default:
		return "LX_" + source.UserID
	}
}
