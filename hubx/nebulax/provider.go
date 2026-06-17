// Package nebulax implements hubx.ClientProvider for "graphx.nebula".
package nebulax

import (
	"context"
	"fmt"

	hubx "github.com/gospacex/hubx"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Endpoint string `yaml:"endpoint"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Provider struct{}

func New() *Provider { return &Provider{} }
func (p *Provider) Name() string { return "graphx.nebula" }

func (p *Provider) Build(instanceName string, cfg map[string]any) (hubx.Client, error) {
	raw, ok := cfg["config"]
	if !ok {
		return nil, fmt.Errorf("%w: graphx.nebula/%s: missing 'config' key", hubx.ErrConfigInvalid, instanceName)
	}
	var c Config
	dec, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml", ErrorUnset: true, ErrorUnused: true, Result: &c,
	})
	if err := dec.Decode(raw); err != nil {
		return nil, fmt.Errorf("%w: graphx.nebula/%s: %v", hubx.ErrConfigInvalid, instanceName, err)
	}
	if c.Endpoint == "" {
		return nil, fmt.Errorf("%w: graphx.nebula/%s: endpoint is required", hubx.ErrConfigInvalid, instanceName)
	}
	return &client{endpoint: c.Endpoint}, nil
}

func (p *Provider) HealthCheck(context.Context) error { return nil }
func (p *Provider) Close() error                      { return nil }

type client struct{ endpoint string }

func (c *client) HealthCheck(ctx context.Context) error { return nil }
func (c *client) Close() error                          { return nil }
