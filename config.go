package graphx

import (
	"crypto/tls"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the common configuration shared by all backends.
type Config struct {
	Name     string         `yaml:"name"`
	Address  string         `yaml:"address"`
	Username string         `yaml:"username"`
	Password string         `yaml:"password"`
	Database string         `yaml:"database"`
	PoolSize int            `yaml:"poolSize"`
	TLS       bool           `yaml:"tls"`
	TLSConfig *tls.Config    `yaml:"-"`
	// TracerName overrides the default OTel tracer name ("graphx/<backend>")
	// used by backend subpackages. Empty means use the backend default.
	TracerName string        `yaml:"tracer_name"`
	Extra     map[string]any `yaml:"extra,omitempty"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	return cfg, yaml.Unmarshal(data, &cfg)
}
