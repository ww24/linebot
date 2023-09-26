package config

import (
	"net/url"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

type ServiceEndpoint struct {
	ServiceEndpoint string   `split_words:"true"`
	serviceEndpoint *url.URL `ignored:"true"`
}

func (c *ServiceEndpoint) Valid() bool {
	return c.ServiceEndpoint != ""
}

func NewServiceEndpoint() (*ServiceEndpoint, error) {
	var conf ServiceEndpoint
	if err := envconfig.Process("", &conf); err != nil {
		return nil, xerrors.Errorf("failed to parse service config: %w", err)
	}
	if conf.ServiceEndpoint != "" {
		u, err := url.Parse(conf.ServiceEndpoint)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse service endpoint: %w", err)
		}
		conf.serviceEndpoint = u
	}
	return &conf, nil
}

func (c *ServiceEndpoint) ResolveServiceEndpoint(path string) (*url.URL, error) {
	if c.serviceEndpoint == nil {
		return nil, xerrors.New("service endpoint is not configured")
	}
	r, err := url.Parse(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse path: %w", err)
	}
	return c.serviceEndpoint.ResolveReference(r), nil
}
