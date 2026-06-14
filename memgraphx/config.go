package memgraphx

import (
	"fmt"
	"strings"

	"github.com/gospacex/graphx"
)

// MemgraphConfig holds the resolved Bolt connection parameters.
type MemgraphConfig struct {
	BoltURI string
}

// Build resolves a graphx.Config into a MemgraphConfig with a bolt:// URI.
func Build(cfg graphx.Config) (MemgraphConfig, error) {
	if cfg.Address == "" {
		return MemgraphConfig{}, fmt.Errorf("graphx/memgraphx: address must not be empty")
	}
	addr := strings.Replace(cfg.Address, "localhost", "127.0.0.1", 1)
	return MemgraphConfig{
		BoltURI: "bolt://" + addr,
	}, nil
}
