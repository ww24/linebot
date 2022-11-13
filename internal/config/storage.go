package config

import (
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

type Storage struct {
	ImageBucket string `split_words:"true" required:"true"`
}

func NewStorage() (*Storage, error) {
	var conf Storage
	if err := envconfig.Process("STORAGE", &conf); err != nil {
		return nil, xerrors.Errorf("failed to parse storage config: %w", err)
	}
	return &conf, nil
}
