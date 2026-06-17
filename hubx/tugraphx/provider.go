// Package tugraphx implements hubx.ClientProvider for "graphx.tugraph".
package tugraphx

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

func New() *Provider                              { return &Provider{} }
func (p *Provider) Name() string                  { return "graphx.tugraph" }
func (p *Provider) HealthCheck(context.Context) error { return nil }
func (p *Provider) Close() error                  { return nil }

func (p *Provider) Build(instanceName string, cfg map[string]any) (hubx.Client, error) {
	raw, ok := cfg["config"]
	if !ok {
		return nil, fmt.Errorf("%w: graphx.tugraph/%s: missing 'config' key", hubx.ErrConfigInvalid, instanceName)
	}
	var c Config
	dec, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml", ErrorUnused: true, Result: &c,
	})
	if err := dec.Decode(raw); err != nil {
		return nil, fmt.Errorf("%w: graphx.tugraph/%s: %v", hubx.ErrConfigInvalid, instanceName, err)
	}
	if c.Endpoint == "" {
		return nil, fmt.Errorf("%w: graphx.tugraph/%s: endpoint is required", hubx.ErrConfigInvalid, instanceName)
	}
	return &client{endpoint: c.Endpoint}, nil
}

type client struct{ endpoint string }

func (c *client) HealthCheck(ctx context.Context) error { return nil }
func (c *client) Close() error                          { return nil }
