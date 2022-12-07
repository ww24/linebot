package config

import (
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

type AccessLog struct {
	Topic          string `split_words:"true"`
	TrustedProxies int    `split_words:"true" default:"0"`
}

func NewAccessLog() (*AccessLog, error) {
	var conf AccessLog
	if err := envconfig.Process("ACCESS_LOG", &conf); err != nil {
		return nil, xerrors.Errorf("failed to parse access log config: %w", err)
	}
	return &conf, nil
}
