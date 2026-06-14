package janusx

import (
	"fmt"
	"strings"

	"github.com/gospacex/graphx"
)

// JanusConfig holds the resolved Gremlin WebSocket URL.
type JanusConfig struct {
	URL string
}

// Build resolves a graphx.Config into a JanusConfig with a ws:// or wss:// URL.
func Build(cfg graphx.Config) (JanusConfig, error) {
	if cfg.Address == "" {
		return JanusConfig{}, fmt.Errorf("graphx/janusx: address must not be empty")
	}
	scheme := "ws"
	if cfg.TLS {
		scheme = "wss"
	}
	addr := strings.Replace(cfg.Address, "localhost", "127.0.0.1", 1)
	return JanusConfig{
		URL: fmt.Sprintf("%s://%s/gremlin", scheme, addr),
	}, nil
}
