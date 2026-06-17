// Package neo4jx implements hubx.ClientProvider for "graphx.neo4j".
//
// The provider wraps github.com/gospacex/graphx/neo4jx as a hubx.Client.
// Build decodes cfg["config"] into graphx.Config via mapstructure
// (TagName: "yaml", ErrorUnset / ErrorUnused both enabled) and calls
// neo4jx.New. Errors are wrapped with the appropriate hubx sentinel:
//   - missing or invalid "config" key → hubx.ErrConfigInvalid
//   - neo4jx.New failure              → hubx.ErrBuildFailed
//
// Provider.HealthCheck and Provider.Close are stateless no-ops because
// the provider owns no shared connection state — every instance gets
// its own *neo4jx.Neo4j handle in Build.
package neo4jx

import (
	"context"
	"fmt"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/neo4jx"
	hubx "github.com/gospacex/hubx"
	"github.com/mitchellh/mapstructure"
)

// Provider implements hubx.ClientProvider for the "graphx.neo4j" driver.
type Provider struct{}

// New returns a new graphx.neo4j Provider.
func New() *Provider { return &Provider{} }

// Name returns the registry name.
func (p *Provider) Name() string { return "graphx.neo4j" }

// Build decodes cfg["config"] into graphx.Config via mapstructure
// (TagName: "yaml", ErrorUnset / ErrorUnused both enabled) and then
// calls neo4jx.New.
//
// Errors are wrapped with the appropriate hubx sentinel:
//   - missing or invalid "config" key → hubx.ErrConfigInvalid
//   - neo4jx.New failure              → hubx.ErrBuildFailed
func (p *Provider) Build(instanceName string, cfg map[string]any) (hubx.Client, error) {
	raw, ok := cfg["config"]
	if !ok {
		return nil, fmt.Errorf("%w: graphx.neo4j/%s: missing 'config' key", hubx.ErrConfigInvalid, instanceName)
	}

	var gc graphx.Config
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:     "yaml",
		ErrorUnused: true,
		Result:      &gc,
		ZeroFields:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: graphx.neo4j/%s: decoder: %v", hubx.ErrConfigInvalid, instanceName, err)
	}
	if err := dec.Decode(raw); err != nil {
		return nil, fmt.Errorf("%w: graphx.neo4j/%s: %v", hubx.ErrConfigInvalid, instanceName, err)
	}

	cli, err := neo4jx.New(context.Background(), gc)
	if err != nil {
		return nil, fmt.Errorf("%w: graphx.neo4j/%s: %v", hubx.ErrBuildFailed, instanceName, err)
	}
	return &client{n: cli}, nil
}

// HealthCheck is a no-op for the provider itself — the provider owns
// no connection state.
func (p *Provider) HealthCheck(context.Context) error { return nil }

// Close is a no-op for the provider itself.
func (p *Provider) Close() error { return nil }

// client wraps a *neo4jx.Neo4j as a hubx.Client.
type client struct{ n *neo4jx.Neo4j }

// HealthCheck delegates to the driver's Ping.
func (c *client) HealthCheck(ctx context.Context) error { return c.n.Ping(ctx) }

// Close delegates to the driver's Close.
func (c *client) Close() error { return c.n.Close(context.Background()) }