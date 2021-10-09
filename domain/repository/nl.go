package repository

import "github.com/ww24/linebot/domain/model"

type NLParser interface {
	Parse(string) *model.Item
}
