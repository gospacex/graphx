package neo4jx

import (
	"fmt"
	"strings"

	"github.com/gospacex/graphx"
)

// Neo4jConfig holds neo4j-specific configuration derived from graphx.Config.
type Neo4jConfig struct {
	BoltURI string
}

// Build returns the bolt URI for the given graphx.Config.
// Scheme defaults to "bolt" unless cfg.Extra["scheme"] is set.
func Build(cfg graphx.Config) (Neo4jConfig, error) {
	if cfg.Address == "" {
		return Neo4jConfig{}, fmt.Errorf("graphx/neo4jx: address must not be empty")
	}
	scheme := "bolt"
	if s, ok := cfg.Extra["scheme"]; ok {
		if str, ok2 := s.(string); ok2 {
			scheme = str
		}
	}
	addr := strings.TrimSpace(cfg.Address)
	return Neo4jConfig{
		BoltURI: fmt.Sprintf("%s://%s", scheme, addr),
	}, nil
}
