package interactor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/domain/service"
	"github.com/ww24/linebot/usecase"
)

const (
	triggerShopping = "買い物リスト"
	prefixShopping  = "【買い物リスト】"
)

type Shopping struct {
	shopping service.Shopping
	nlParser repository.NLParser
	message  repository.MessageProviderSet
	bot      service.Bot
}

func NewShopping(
	shopping service.Shopping,
	nlParser repository.NLParser,
	message repository.MessageProviderSet,
	bot service.Bot,
) *Shopping {
	return &Shopping{
		shopping: shopping,
		nlParser: nlParser,
		message:  message,
		bot:      bot,
	}
}

func (s *Shopping) Handle(ctx context.Context, e *model.Event) error {
	err := e.HandleTypeMessage(ctx, func(ctx context.Context, e *model.Event) error {
		if e.FilterText(triggerShopping) {
			return s.handleMenu(ctx, e)
		}

		return s.handleStatus(ctx, e)
	})
	if err != nil {
		return err
	}

	if err := e.HandleTypePostback(ctx, s.handlePostBack); err != nil {
		return err
	}

	return nil
}

func (s *Shopping) handleMenu(ctx context.Context, e *model.Event, texts ...string) error {
	items, err := s.shopping.List(ctx, e.ConversationID())
	if err != nil {
		return xerrors.Errorf("failed to list shopping items: %w", err)
	}

	prefixMsg := prefixShopping
	if len(texts) > 0 {
		prefixMsg += strings.Join(texts, "\n") + "\n\n"
	}

	if len(items) == 0 {
		text := prefixMsg + "リストは空です。\n何をしますか？"
		msg := s.message.ShoppingMenu(text, model.ShoppingReplyTypeEmptyList)
		if err := s.bot.ReplyMessage(ctx, e, msg); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
		return nil
	}

	text := fmt.Sprintf(prefixMsg+"%d件登録されています。\n%s\n\n何をしますか？",
		len(items), model.ShoppingItems(items).Print(model.ListTypeOrdered))
	msg := s.message.ShoppingMenu(text, model.ShoppingReplyTypeWithoutView)
	if err := s.bot.ReplyMessage(ctx, e, msg); err != nil {
		return xerrors.Errorf("failed to reply message: %w", err)
	}

	return nil
}

func (s *Shopping) handlePostBack(ctx context.Context, e *model.Event) error {
	conversationID := e.ConversationID()

	switch e.Postback.Data {
	case "Shopping#delete":
		text := prefixShopping + "リストを空にしても良いですか？"
		msg := s.message.ShoppingDeleteConfirmation(text)
		if err := s.bot.ReplyMessage(ctx, e, msg); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}

	case "Shopping#deleteConfirm":
		if err := s.shopping.DeleteAllItem(ctx, conversationID); err != nil {
			return xerrors.Errorf("failed to delete all shopping items: %w", err)
		}
		if err := s.handleMenu(ctx, e); err != nil {
			return err
		}

	case "Shopping#deleteCancel":
		if err := s.shopping.SetStatusShopping(ctx, conversationID); err != nil {
			return xerrors.Errorf("failed to set status: %w", err)
		}
		if err := s.handleMenu(ctx, e); err != nil {
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
		if err := s.bot.ReplyTextMessage(ctx, e, text); err != nil {
			return xerrors.Errorf("failed to reply text message: %w", err)
		}

	case "Shopping#view":
		items, err := s.shopping.List(ctx, conversationID)
		if err != nil {
			return xerrors.Errorf("failed to list shopping items: %w", err)
		}

		text := prefixShopping + "\n" + model.ShoppingItems(items).Print(model.ListTypeOrdered)
		msg := s.message.ShoppingMenu(text, model.ShoppingReplyTypeWithoutView)
		if err := s.bot.ReplyMessage(ctx, e, msg); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
	}

	return nil
}

func (s *Shopping) handleStatus(ctx context.Context, e *model.Event) error {
	status, err := s.shopping.GetStatus(ctx, e.ConversationID())
	if err != nil {
		return xerrors.Errorf("failed to get status: %w", err)
	}

	switch status.Type {
	case model.ConversationStatusTypeShopping:
		itemText := strings.Join(e.ReadTextLines(), " ")
		// parse message text
		item := s.nlParser.Parse(itemText)
		if err := s.handleMessageAction(ctx, e, item); err != nil {
			return xerrors.Errorf("failed to handle message action: %w", err)
		}

	case model.ConversationStatusTypeShoppingAdd:
		lines := e.ReadTextLines()
		items := make([]*model.ShoppingItem, 0, len(lines))
		for i, line := range lines {
			item := &model.ShoppingItem{
				ConversationID: e.ConversationID(),
				Name:           line,
				CreatedAt:      time.Now().Unix(),
				Order:          i,
			}
			items = append(items, item)
		}
		if err := s.shopping.AddItem(ctx, e.ConversationID(), items...); err != nil {
			return xerrors.Errorf("failed to add item: %w", err)
		}

		text := fmt.Sprintf(prefixShopping+"%d件追加されました。", len(lines))
		msg := s.message.ShoppingMenu(text, model.ShoppingReplyTypeAll)
		if err := s.bot.ReplyMessage(ctx, e, msg); err != nil {
			return xerrors.Errorf("failed to reply message: %w", err)
		}
	}

	return nil
}

func (s *Shopping) handleMessageAction(ctx context.Context, e *model.Event, item *model.Item) error {
	switch item.Action {
	case model.ActionTypeDelete:
		foundItems, err := s.deleteFromItem(ctx, e.ConversationID(), item)
		if err != nil {
			return err
		}
		text := "次の商品を削除しました。\n" + foundItems.Print(model.ListTypeDotted)
		if err := s.handleMenu(ctx, e, text); err != nil {
			return err
		}

		return nil

	default:
		// do nothing
		return nil
	}
}

func (s *Shopping) deleteFromItem(ctx context.Context, conversationID model.ConversationID, item *model.Item) (model.ShoppingItems, error) {
	items, err := s.shopping.List(ctx, conversationID)
	if err != nil {
		return nil, xerrors.Errorf("failed to list shopping items: %w", err)
	}

	ret := make([]*model.ShoppingItem, 0)

	indexes := item.UniqueIndexes()
	if len(indexes) == 0 {
		return ret, xerrors.Errorf("item not found: %w", usecase.ErrItemNotFound)
	}

	ids := make([]string, 0, len(indexes))
	for _, idx := range indexes {
		if idx <= 0 || idx > len(items) {
			continue
		}
		item := items[idx-1]
		ret = append(ret, item)
		ids = append(ids, item.ID)
	}
	if len(ids) == 0 {
		return ret, nil
	}
	if err := s.shopping.DeleteItems(ctx, conversationID, ids); err != nil {
		return nil, xerrors.Errorf("failed to delete shopping items: %w", err)
	}

	return ret, nil
}
