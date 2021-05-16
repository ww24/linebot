package bot

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/nl"
	"golang.org/x/xerrors"
)

const (
	triggerShopping = "買い物リスト"
	prefixShopping  = "【買い物リスト】"
)

var (
	errNotFound = errors.New("shopping item not found")
)

type ShoppingService struct {
	shopping service.Shopping
	nlParser *nl.Parser
	now      func() time.Time
}

func NewShoppingService(shopping service.Shopping) (*ShoppingService, error) {
	parser, err := nl.NewParser()
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize nl parser: %w", err)
	}
	return &ShoppingService{
		shopping: shopping,
		nlParser: parser,
		now:      time.Now,
	}, nil
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

func (s *ShoppingService) handleMenu(ctx context.Context, bot *Bot, e *linebot.Event, texts ...string) error {
	items, err := s.shopping.List(ctx, ConversationID(e.Source))
	if err != nil {
		return xerrors.Errorf("failed to list shopping items: %w", err)
	}

	prefixMsg := prefixShopping
	if len(texts) > 0 {
		prefixMsg += strings.Join(texts, "\n") + "\n\n"
	}

	if len(items) == 0 {
		var msg linebot.SendingMessage
		msg = linebot.NewTextMessage(prefixMsg + "リストは空です。\n何をしますか？")
		msg = s.addQuickReplies(msg, shoppingRepliesTypeEmptyList)
		c := bot.cli.ReplyMessage(e.ReplyToken, msg)
		if _, err := c.WithContext(ctx).Do(); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
		return nil
	}

	var msg linebot.SendingMessage
	text := fmt.Sprintf(prefixMsg+"%d件登録されています。\n%s\n\n何をしますか？",
		len(items), model.ShoppingItems(items).Print(model.ListTypeOrdered))
	msg = linebot.NewTextMessage(text)
	msg = s.addQuickReplies(msg, shoppingRepliesTypeWithoutView)
	c := bot.cli.ReplyMessage(e.ReplyToken, msg)
	if _, err := c.WithContext(ctx).Do(); err != nil {
		return xerrors.Errorf("failed to reply message: %w", err)
	}

	return nil
}

func (s *ShoppingService) handlePostBack(ctx context.Context, bot *Bot, e *linebot.Event) error {
	conversationID := ConversationID(e.Source)

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
			return xerrors.Errorf("failed to reply message: %w", err)
		}

	case "Shopping#deleteConfirm":
		if err := s.shopping.DeleteAllItem(ctx, conversationID); err != nil {
			return xerrors.Errorf("failed to delete all shopping items: %w", err)
		}
		if err := s.handleMenu(ctx, bot, e); err != nil {
			return err
		}

	case "Shopping#deleteCancel":
		if err := s.shopping.SetStatusShopping(ctx, conversationID); err != nil {
			return xerrors.Errorf("failed to set status: %w", err)
		}
		if err := s.handleMenu(ctx, bot, e); err != nil {
			return err
		}

	case "Shopping#add":
		status := &model.ConversationStatus{
			ConversationID: conversationID,
			Type:           model.ConversationStatusTypeShoppingAdd,
		}
		if err := s.shopping.SetStatus(ctx, status); err != nil {
			return xerrors.Errorf("failed to set status: %w", err)
		}
		text := prefixShopping + "追加する商品を1行に1つずつ入力してください。"
		if err := bot.replyTestMessage(ctx, e, text); err != nil {
			return err
		}

	case "Shopping#view":
		items, err := s.shopping.List(ctx, conversationID)
		if err != nil {
			return xerrors.Errorf("failed to list shopping items: %w", err)
		}

		var msg linebot.SendingMessage
		text := prefixShopping + "\n" + model.ShoppingItems(items).Print(model.ListTypeOrdered)
		msg = linebot.NewTextMessage(text)
		msg = s.addQuickReplies(msg, shoppingRepliesTypeWithoutView)
		c := bot.cli.ReplyMessage(e.ReplyToken, msg)
		if _, err := c.WithContext(ctx).Do(); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
	}

	return nil
}

func (s *ShoppingService) handleStatus(ctx context.Context, bot *Bot, e *linebot.Event) error {
	conversationID := ConversationID(e.Source)
	status, err := s.shopping.GetStatus(ctx, conversationID)
	if err != nil {
		return xerrors.Errorf("failed to get status: %w", err)
	}

	switch status.Type {
	case model.ConversationStatusTypeShopping:
		itemText := strings.Join(bot.readTextLines(e), " ")
		item := s.nlParser.Parse(itemText)
		if item.Action != nl.ActionTypeDelete {
			return nil
		}
		foundItems, err := s.deleteFromItem(ctx, conversationID, item)
		if err != nil {
			return err
		}
		text := "次の商品を削除しました。\n" + foundItems.Print(model.ListTypeDotted)
		if err := s.handleMenu(ctx, bot, e, text); err != nil {
			return err
		}

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
		if err := s.shopping.AddItem(ctx, conversationID, items...); err != nil {
			return xerrors.Errorf("failed to add item: %w", err)
		}

		var msg linebot.SendingMessage
		text := fmt.Sprintf(prefixShopping+"%d件追加されました。", len(lines))
		msg = linebot.NewTextMessage(text)
		msg = s.addQuickReplies(msg, shoppingRepliesTypeAll)
		c := bot.cli.ReplyMessage(e.ReplyToken, msg)
		if _, err := c.WithContext(ctx).Do(); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
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

func (s *ShoppingService) deleteFromItem(ctx context.Context, conversationID model.ConversationID, item *nl.Item) (model.ShoppingItems, error) {
	items, err := s.shopping.List(ctx, conversationID)
	if err != nil {
		return nil, xerrors.Errorf("failed to list shopping items: %w", err)
	}

	ret := make([]*model.ShoppingItem, 0)

	if len(item.Indexes) > 0 {
		ids := make([]string, 0, len(item.Indexes))
		for _, idx := range item.Indexes {
			if idx <= 0 || idx > len(items) {
				continue
			}
			item := items[idx-1]
			ret = append(ret, item)
			ids = append(ids, item.ID)
		}
		if err := s.shopping.DeleteItems(ctx, conversationID, ids); err != nil {
			return nil, xerrors.Errorf("failed to delete shopping items: %w", err)
		}

		return ret, nil
	}

	// TODO: search by name
	// 固有名詞が分割されてしまうので実装が難しい

	return nil, errNotFound
}
