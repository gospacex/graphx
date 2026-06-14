package nebulax

import (
	"fmt"
	"net"
	"strings"

	"github.com/gospacex/graphx"
)

// NebulaConfig holds the resolved host list and graph space name.
type NebulaConfig struct {
	Hosts []string
	Space string
}

// Build resolves a graphx.Config into a NebulaConfig with host:port pairs and a space name.
func Build(cfg graphx.Config) (NebulaConfig, error) {
	if cfg.Address == "" {
		return NebulaConfig{}, fmt.Errorf("graphx/nebulax: address must not be empty")
	}

	space := cfg.Database
	if space == "" {
		space = "default"
	}

	addrs := strings.Split(cfg.Address, ",")
	var hosts []string
	for _, a := range addrs {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		a = strings.Replace(a, "localhost", "127.0.0.1", 1)
		host, port, err := net.SplitHostPort(a)
		if err != nil {
			host = a
			port = "9669"
		}
		_ = host
		hosts = append(hosts, net.JoinHostPort(host, port))
	}
	if len(hosts) == 0 {
		hosts = []string{"127.0.0.1:9669"}
	}

	return NebulaConfig{
		Hosts: hosts,
		Space: space,
	}, nil
}
