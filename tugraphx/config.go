package tugraphx

import (
	"fmt"
	"strings"

	"github.com/gospacex/graphx"
)

// TuGraphConfig holds the resolved Bolt connection parameters.
type TuGraphConfig struct {
	BoltURI string
}

// Build resolves a graphx.Config into a TuGraphConfig.
func Build(cfg graphx.Config) (TuGraphConfig, error) {
	if cfg.Address == "" {
		return TuGraphConfig{}, fmt.Errorf("graphx/tugraphx: address must not be empty")
	}

	addr := cfg.Address
	if strings.HasPrefix(addr, "localhost") {
		addr = strings.Replace(addr, "localhost", "127.0.0.1", 1)
	}

	scheme := "bolt"
	if cfg.TLS {
		scheme = "bolts"
	}

	return TuGraphConfig{
		BoltURI: fmt.Sprintf("%s://%s", scheme, addr),
	}, nil
}
