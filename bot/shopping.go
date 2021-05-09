package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
)

const (
	triggerShopping = "買い物リスト"
	prefixShopping  = "【買い物リスト】"
)

type ShoppingService struct {
	conversation repository.Conversation
	now          func() time.Time
}

func NewShoppingService(conversation repository.Conversation) *ShoppingService {
	return &ShoppingService{
		conversation: conversation,
		now:          time.Now,
	}
}

func (s *ShoppingService) Handle(ctx context.Context, bot *Bot, e *linebot.Event) error {
	switch e.Type {
	case linebot.EventTypeMessage:
		if bot.filterText(e, triggerShopping) {
			if err := s.handleMenu(ctx, bot, e); err != nil {
				return err
			}
		} else {
			if err := s.handleStatus(ctx, bot, e); err != nil {
				return err
			}
		}

	case linebot.EventTypePostback:
		if err := s.handlePostBack(ctx, bot, e); err != nil {
			return err
		}
	}

	return nil
}

func (s *ShoppingService) handleMenu(ctx context.Context, bot *Bot, e *linebot.Event) error {
	if err := s.setStatusShopping(ctx, bot, e); err != nil {
		return err
	}

	items, err := s.conversation.FindShoppingItem(ctx, ConversationID(e.Source))
	if err != nil {
		return err
	}

	if len(items) == 0 {
		var msg linebot.SendingMessage
		msg = linebot.NewTextMessage(prefixShopping + "リストは空です。\n何をしますか？")
		msg = s.addQuickReplies(msg, shoppingRepliesTypeEmptyList)
		c := bot.cli.ReplyMessage(e.ReplyToken, msg)
		if _, err := c.WithContext(ctx).Do(); err != nil {
			return err
		}
		return nil
	}

	var msg linebot.SendingMessage
	text := fmt.Sprintf(prefixShopping+"%d件登録されています。\n%s\n\n何をしますか？",
		len(items), model.ShoppingItems(items).Print())
	msg = linebot.NewTextMessage(text)
	msg = s.addQuickReplies(msg, shoppingRepliesTypeWithoutView)
	c := bot.cli.ReplyMessage(e.ReplyToken, msg)
	if _, err := c.WithContext(ctx).Do(); err != nil {
		return err
	}

	return nil
}

func (s *ShoppingService) handlePostBack(ctx context.Context, bot *Bot, e *linebot.Event) error {
	switch e.Postback.Data {
	case "Shopping#delete":
		var msg linebot.SendingMessage
		text := prefixShopping + "リストを空にしても良いですか？"
		msg = linebot.NewTextMessage(text)
		msg = msg.WithQuickReplies(&linebot.QuickReplyItems{
			Items: []*linebot.QuickReplyButton{
				{Action: linebot.NewPostbackAction("YES", "Shopping#deleteConfirm", "", "YES")},
				{Action: linebot.NewPostbackAction("NO", "Shopping#deleteCancel", "", "NO")},
			},
		})
		c := bot.cli.ReplyMessage(e.ReplyToken, msg)
		if _, err := c.WithContext(ctx).Do(); err != nil {
			return err
		}

	case "Shopping#deleteConfirm":
		if err := s.conversation.DeleteAllShoppingItem(ctx, ConversationID(e.Source)); err != nil {
			return err
		}
		if err := s.setStatusShopping(ctx, bot, e); err != nil {
			return err
		}
		if err := s.handleMenu(ctx, bot, e); err != nil {
			return err
		}

	case "Shopping#deleteCancel":
		if err := s.setStatusShopping(ctx, bot, e); err != nil {
			return err
		}
		if err := s.handleMenu(ctx, bot, e); err != nil {
			return err
		}

	case "Shopping#add":
		status := &model.ConversationStatus{
			ConversationID: ConversationID(e.Source),
			Type:           model.ConversationStatusTypeShoppingAdd,
		}
		if err := s.conversation.SetStatus(ctx, status); err != nil {
			return err
		}
		text := prefixShopping + "追加する商品を1行に1つずつ入力してください。"
		if err := bot.replyTestMessage(ctx, e, text); err != nil {
			return err
		}

	case "Shopping#view":
		items, err := s.conversation.FindShoppingItem(ctx, ConversationID(e.Source))
		if err != nil {
			return err
		}

		var msg linebot.SendingMessage
		text := prefixShopping + "\n" + model.ShoppingItems(items).Print()
		msg = linebot.NewTextMessage(text)
		msg = s.addQuickReplies(msg, shoppingRepliesTypeWithoutView)
		c := bot.cli.ReplyMessage(e.ReplyToken, msg)
		if _, err := c.WithContext(ctx).Do(); err != nil {
			return err
		}
	}

	return nil
}

func (s *ShoppingService) handleStatus(ctx context.Context, bot *Bot, e *linebot.Event) error {
	status, err := s.conversation.GetStatus(ctx, ConversationID(e.Source))
	if err != nil {
		return err
	}

	switch status.Type {
	case model.ConversationStatusTypeShoppingAdd:
		lines := bot.readTextLines(e)
		items := make([]*model.ShoppingItem, 0, len(lines))
		for _, line := range lines {
			item := &model.ShoppingItem{
				ConversationID: ConversationID(e.Source),
				Name:           line,
				CreatedAt:      s.now().Unix(),
			}
			items = append(items, item)
		}
		if err := s.conversation.AddShoppingItem(ctx, items...); err != nil {
			return err
		}
		if err := s.setStatusShopping(ctx, bot, e); err != nil {
			return err
		}

		var msg linebot.SendingMessage
		text := fmt.Sprintf(prefixShopping+"%d件追加されました。", len(lines))
		msg = linebot.NewTextMessage(text)
		msg = s.addQuickReplies(msg, shoppingRepliesTypeAll)
		c := bot.cli.ReplyMessage(e.ReplyToken, msg)
		if _, err := c.WithContext(ctx).Do(); err != nil {
			return err
		}
	}

	return nil
}

func (s *ShoppingService) setStatusShopping(ctx context.Context, bot *Bot, e *linebot.Event) error {
	status := &model.ConversationStatus{
		ConversationID: ConversationID(e.Source),
		Type:           model.ConversationStatusTypeShopping,
	}
	if err := s.conversation.SetStatus(ctx, status); err != nil {
		return err
	}
	return nil
}

type shoppingRepliesType int

const (
	shoppingRepliesTypeAll shoppingRepliesType = iota
	shoppingRepliesTypeEmptyList
	shoppingRepliesTypeWithoutView
)

func (s *ShoppingService) addQuickReplies(msg linebot.SendingMessage, typ shoppingRepliesType) linebot.SendingMessage {
	var items []*linebot.QuickReplyButton
	switch typ {
	case shoppingRepliesTypeEmptyList:
		items = []*linebot.QuickReplyButton{
			{Action: linebot.NewPostbackAction("追加", "Shopping#add", "", "追加")},
		}
	case shoppingRepliesTypeWithoutView:
		items = []*linebot.QuickReplyButton{
			{Action: linebot.NewPostbackAction("削除", "Shopping#delete", "", "削除")},
			{Action: linebot.NewPostbackAction("追加", "Shopping#add", "", "追加")},
		}
	default:
		items = []*linebot.QuickReplyButton{
			{Action: linebot.NewPostbackAction("削除", "Shopping#delete", "", "削除")},
			{Action: linebot.NewPostbackAction("追加", "Shopping#add", "", "追加")},
			{Action: linebot.NewPostbackAction("表示", "Shopping#view", "", "表示")},
		}
	}

	return msg.WithQuickReplies(&linebot.QuickReplyItems{
		Items: items,
	})
}
