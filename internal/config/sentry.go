package config

import (
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

type Sentry struct {
	Enable bool   `split_words:"true" default:"false"`
	DSN    string `split_words:"true"`
}

func NewSentry() (*Sentry, error) {
	var conf Sentry
	if err := envconfig.Process("SENTRY", &conf); err != nil {
		return nil, xerrors.Errorf("failed to parse sentry config: %w", err)
	}
	return &conf, nil
}
