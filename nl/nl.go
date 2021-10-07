package nl

import (
	"strconv"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"golang.org/x/text/unicode/norm"
)

const (
	// POS 1
	posNoun = "名詞"
	posVerb = "動詞"

	// POS 2
	posNumeral = "数"  // 名詞, 数
	posSuffix  = "接尾" // 名詞, 接尾

	// POS 3
	posQuantifier = "助数詞" // 名詞, 接尾, 助数詞
)

type ActionType int

const (
	ActionTypeUnknown ActionType = iota
	ActionTypeDelete
)

type Item struct {
	Indexes []int
	Name    []string
	Action  ActionType
}

type Parser struct {
	tokenizer *tokenizer.Tokenizer
	allow     *filter.POSFilter
	deny      *filter.POSFilter
	replacer  *strings.Replacer
}

func NewParser() (*Parser, error) {
	tk, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}

	allowFilter := filter.NewPOSFilter(
		filter.POS{posNoun},
		filter.POS{posVerb},
	)
	denyFilter := filter.NewPOSFilter(
		filter.POS{posNoun, posSuffix, posQuantifier},
	)
	replacer := strings.NewReplacer(",", "、")

	return &Parser{
		tokenizer: tk,
		allow:     allowFilter,
		deny:      denyFilter,
		replacer:  replacer,
	}, nil
}

func (p *Parser) Parse(str string) *Item {
	str = norm.NFKC.String(str)
	str = p.replacer.Replace(str)
	tokens := p.tokenizer.Tokenize(str)

	// debug code
	// for _, token := range tokens {
	// 	fmt.Printf("%+v, %+v\n", token.Surface, token.Features())
	// }

	p.allow.Keep(&tokens)
	p.deny.Drop(&tokens)

	item := &Item{}
	for _, token := range tokens {
		at := p.selectAction(token)
		if at != ActionTypeUnknown {
			item.Action = at
			continue
		}

		pos := token.POS()
		switch pos[0] {
		case posNoun:
			if pos[1] == posNumeral {
				idx, err := p.parseNumber(token.Surface)
				if err != nil {
					continue
				}
				item.Indexes = append(item.Indexes, idx)
				continue
			}

			// add name if reading feature exists
			if _, ok := token.Reading(); ok {
				item.Name = append(item.Name, token.Surface)
			}
		}
	}

	return item
}

func (*Parser) parseNumber(str string) (int, error) {
	// 漢数字対応
	num, err := strconv.Atoi(str)
	if err != nil {
		for i, n := range cn {
			if str == n {
				return i + 1, nil
			}
		}

		return 0, err
	}
	return num, nil
}

func (*Parser) selectAction(t tokenizer.Token) ActionType {
	keyword := ""
	if bf, ok := t.BaseForm(); ok {
		keyword = bf
	} else {
		keyword = t.Surface
	}

	switch keyword {
	case "削除", "除去", "消す":
		return ActionTypeDelete
	default:
		return ActionTypeUnknown
	}
}

var cn = [...]string{"一", "二", "三", "四", "五", "六", "七", "八", "九"}
