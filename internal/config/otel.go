package config

import (
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

type Otel struct {
	SamplingRate float64 `split_words:"true" default:"0.1"`
}

func NewOtel() (*Otel, error) {
	var conf Otel
	if err := envconfig.Process("OTEL", &conf); err != nil {
		return nil, xerrors.Errorf("failed to parse otel config: %w", err)
	}
	return &conf, nil
}
