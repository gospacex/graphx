package agex

import (
	"fmt"
	"net/url"

	"github.com/gospacex/graphx"
)

// AGEConfig holds the resolved PostgreSQL DSN and graph name.
type AGEConfig struct {
	DSN   string
	Graph string
}

// Build resolves a graphx.Config into an AGEConfig with a PostgreSQL DSN.
func Build(cfg graphx.Config) (AGEConfig, error) {
	if cfg.Address == "" {
		return AGEConfig{}, fmt.Errorf("graphx/agex: address must not be empty")
	}

	password := url.QueryEscape(cfg.Password)
	sslmode := "disable"
	if cfg.TLS {
		sslmode = "require"
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Username, password, cfg.Address, cfg.Database, sslmode)

	graph := "graphx"
	if g, ok := cfg.Extra["graph"]; ok {
		if s, ok2 := g.(string); ok2 && s != "" {
			graph = s
		}
	}

	return AGEConfig{DSN: dsn, Graph: graph}, nil
}
