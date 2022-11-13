package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

const defaultTimezoneOffset = 9 * 60 * 60

//nolint:gochecknoglobals
var defaultLocation = time.FixedZone("Asia/Tokyo", defaultTimezoneOffset)

type Time struct {
	DefaultTimezone string `split_words:"true"`
}

func NewTime() (*Time, error) {
	var conf Time
	if err := envconfig.Process("TIME", &conf); err != nil {
		return nil, xerrors.Errorf("failed to parse time config: %w", err)
	}
	return &conf, nil
}

func (c *Time) DefaultLocation() *time.Location {
	if c.DefaultTimezone == "" {
		return defaultLocation
	}
	loc, err := time.LoadLocation(c.DefaultTimezone)
	if err != nil {
		return defaultLocation
	}
	return loc
}
