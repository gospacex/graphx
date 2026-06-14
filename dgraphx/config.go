package dgraphx

import (
	"fmt"
	"strings"

	"github.com/gospacex/graphx"
)

// DgraphConfig holds the resolved HTTP endpoint parameters.
type DgraphConfig struct {
	BaseURL string
	Scheme  string
}

// Build resolves a graphx.Config into a DgraphConfig with an HTTP base URL.
// Port 9080 (gRPC) is auto-rewritten to 8080 (HTTP).
func Build(cfg graphx.Config) (DgraphConfig, error) {
	if cfg.Address == "" {
		return DgraphConfig{}, fmt.Errorf("graphx/dgraphx: address must not be empty")
	}
	addr := strings.Replace(cfg.Address, "localhost", "127.0.0.1", 1)

	if strings.HasSuffix(addr, ":9080") {
		addr = strings.Replace(addr, ":9080", ":8080", 1)
	}

	scheme := "http"
	if cfg.TLS {
		scheme = "https"
	}

	return DgraphConfig{
		BaseURL: fmt.Sprintf("%s://%s", scheme, addr),
		Scheme:  scheme,
	}, nil
}
