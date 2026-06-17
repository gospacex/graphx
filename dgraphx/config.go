package dgraphx

import (
	"fmt"
	"net"

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

	host, port, err := net.SplitHostPort(cfg.Address)
	if err != nil {
		host = cfg.Address
		port = "8080"
	}

	if port == "9080" {
		port = "8080"
	}

	if host == "localhost" {
		host = "127.0.0.1"
	}

	scheme := "http"
	if cfg.TLS {
		scheme = "https"
	}

	addr := net.JoinHostPort(host, port)
	return DgraphConfig{
		BaseURL: fmt.Sprintf("%s://%s", scheme, addr),
		Scheme:  scheme,
	}, nil
}
