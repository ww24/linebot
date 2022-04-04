package linebot

import (
	_ "embed"
	"encoding/json"
	"time"

	"github.com/google/go-jsonnet"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
)

var (
	//go:embed flex/reminder_list.jsonnet
	reminderListMessage string
)

type ReminderItem struct {
	Title        string `json:"title"`
	SubTitle     string `json:"subTitle"`
	Next         string `json:"next"`
	DeleteTarget string `json:"deleteTarget"`
}

func makeReminderListMessage(items []*model.ReminderItem, t time.Time) ([]byte, error) {
	reminderItems := make([]*ReminderItem, 0, len(items))
	for _, item := range items {
		reminderItems = append(reminderItems, toReminderItem(item, t))
	}

	itemsJSON, err := json.Marshal(reminderItems)
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal reminder items: %w", err)
	}

	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.MemoryImporter{
		Data: map[string]jsonnet.Contents{
			"reminder_list.json": jsonnet.MakeContents(string(itemsJSON)),
		},
	})

	result, err := vm.EvaluateAnonymousSnippet("reminder_list.jsonnet", reminderListMessage)
	if err != nil {
		return nil, xerrors.Errorf("failed to evaluate jsonnet: %w", err)
	}

	return []byte(result), nil
}

func toReminderItem(item *model.ReminderItem, t time.Time) *ReminderItem {
	var next string
	if schedule, err := item.Scheduler.Next(t); err != nil {
		next = "ERROR: failed to calculate next schedule"
	} else {
		next = schedule.Format("01/02 15:04")
	}
	return &ReminderItem{
		Title:        item.Executor.Type.UIText(),
		SubTitle:     item.Scheduler.UIText(),
		Next:         next,
		DeleteTarget: "Reminder#delete#" + string(item.ID),
	}
}
