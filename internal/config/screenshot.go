package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

type Screenshot struct {
	TargetURL      string        `split_words:"true" required:"true"`
	TargetSelector string        `split_words:"true" required:"true"`
	BrowserTimeout time.Duration `split_words:"true" default:"60s"`
}

func NewScreenshot() (*Screenshot, error) {
	var conf Screenshot
	if err := envconfig.Process("SCREENSHOT", &conf); err != nil {
		return nil, xerrors.Errorf("failed to parse screenshot config: %w", err)
	}
	return &conf, nil
}
