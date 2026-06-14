package tigergraphx

import (
	"fmt"
	"strings"

	"github.com/gospacex/graphx"
)

// TigerGraphConfig holds the resolved HTTP base URL.
type TigerGraphConfig struct {
	BaseURL string
}

// Build resolves a graphx.Config into a TigerGraphConfig with an http:// or https:// base URL.
func Build(cfg graphx.Config) (TigerGraphConfig, error) {
	if cfg.Address == "" {
		return TigerGraphConfig{}, fmt.Errorf("graphx/tigergraphx: address must not be empty")
	}
	scheme := "http"
	if cfg.TLS {
		scheme = "https"
	}
	addr := strings.Replace(cfg.Address, "localhost", "127.0.0.1", 1)
	return TigerGraphConfig{
		BaseURL: fmt.Sprintf("%s://%s", scheme, addr),
	}, nil
}
