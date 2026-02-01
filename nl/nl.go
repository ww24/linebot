package nl

import (
	"strconv"
	"strings"

	"github.com/google/wire"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
)

// Set provides a wire set.
var Set = wire.NewSet(
	NewParser,
	wire.Bind(new(repository.NLParser), new(*Parser)),
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

type Parser struct {
	tokenizer *tokenizer.Tokenizer
	allow     *filter.POSFilter
	deny      *filter.POSFilter
	replacer  *strings.Replacer
}

func NewParser() (*Parser, error) {
	tk, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize tokenizer: %w", err)
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

func (p *Parser) Parse(str string) *model.Item {
	str = norm.NFKC.String(str)
	str = p.replacer.Replace(str)
	tokens := p.tokenizer.Tokenize(str)

	// debug code
	// for _, token := range tokens {
	// 	fmt.Printf("%+v, %+v\n", token.Surface, token.Features())
	// }

	p.allow.Keep(&tokens)
	p.deny.Drop(&tokens)

	item := new(model.Item)
	for i := range tokens {
		token := &tokens[i]
		at := p.selectAction(token)
		if at != model.ActionTypeUnknown {
			item.Action = at
			continue
		}

		pos := token.POS()
		if pos[0] == posNoun {
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
	var cn = [...]string{"一", "二", "三", "四", "五", "六", "七", "八", "九"}
	num, err := strconv.Atoi(str)
	if err != nil {
		for i, n := range cn {
			if str == n {
				return i + 1, nil
			}
		}

		return 0, xerrors.Errorf("failed to convert Chinese numeral: %w", err)
	}
	return num, nil
}

func (*Parser) selectAction(t *tokenizer.Token) model.ActionType {
	keyword := ""
	if bf, ok := t.BaseForm(); ok {
		keyword = bf
	} else {
		keyword = t.Surface
	}

	switch keyword {
	case "削除", "除去", "消す":
		return model.ActionTypeDelete
	default:
		return model.ActionTypeUnknown
	}
}
