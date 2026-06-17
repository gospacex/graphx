package tigergraphx

import (
	"fmt"
	"net"

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

	host, port, err := net.SplitHostPort(cfg.Address)
	if err != nil {
		host = cfg.Address
		port = "14240"
	}

	if host == "localhost" {
		host = "127.0.0.1"
	}

	scheme := "http"
	if cfg.TLS {
		scheme = "https"
	}

	return TigerGraphConfig{
		BaseURL: fmt.Sprintf("%s://%s", scheme, net.JoinHostPort(host, port)),
	}, nil
}
