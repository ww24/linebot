//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOPACKAGE/mock_$GOFILE -package=mock_repository

package repository

import "github.com/ww24/linebot/domain/model"

type NLParser interface {
	Parse(string) *model.Item
}
