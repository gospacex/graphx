// Package observability provides OTel tracing integration, span middleware, and exporter factories.
package observability

import "fmt"

// Config holds the observability configuration (exporter type, endpoint, sampler, etc.).
// Can be populated from YAML or built programmatically.
type Config struct {
	Enabled      bool       `yaml:"enabled"`
	Service      string     `yaml:"service_name"`
	Endpoint     string     `yaml:"endpoint"`
	Protocol     string     `yaml:"protocol"`
	Exporter     string     `yaml:"exporter"`

	SamplerType  string     `yaml:"sampler_type"`
	SamplerRatio float64    `yaml:"sampler_ratio"`

	KafkaBrokers []string   `yaml:"kafka_brokers"`
	KafkaTopic   string     `yaml:"kafka_topic"`
	KafkaSASL    SASLConfig `yaml:"kafka_sasl,omitempty"`

	RedisAddrs    []string `yaml:"redis_addrs"`
	RedisStream   string   `yaml:"redis_stream"`
	RedisPassword string   `yaml:"redis_password,omitempty"`
}

// SASLConfig holds Kafka SASL authentication parameters.
type SASLConfig struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Mechanism string `yaml:"mechanism"`
}

// ExporterType identifies the OTLP trace exporter backend.
type ExporterType string

const (
	ExporterJaeger      ExporterType = "jaeger"
	ExporterKafkaTopic  ExporterType = "kafkatopic"
	ExporterRedisStream ExporterType = "redisstream"
)

func (c *Config) setDefaults() {
	if c.Service == "" {
		c.Service = "graphx"
	}
	if c.Endpoint == "" {
		c.Endpoint = "localhost:4318"
	}
	if c.Protocol == "" {
		c.Protocol = "http"
	}
	if c.Exporter == "" {
		c.Exporter = "jaeger"
	}
	if c.SamplerType == "" {
		c.SamplerType = "always_on"
	}
	if c.SamplerRatio <= 0 {
		c.SamplerRatio = 1.0
	}
	if c.KafkaTopic == "" {
		c.KafkaTopic = "otel-traces"
	}
	if c.RedisStream == "" {
		c.RedisStream = "otel-traces"
	}
}

// Validate checks the configuration and sets defaults for empty fields.
func (c *Config) Validate() error {
	c.setDefaults()
	switch c.GetExporterType() {
	case ExporterJaeger, ExporterKafkaTopic, ExporterRedisStream:
	default:
		return fmt.Errorf("observability: unknown exporter %q", c.Exporter)
	}
	return nil
}

// GetExporterType returns the parsed ExporterType from the config string.
func (c *Config) GetExporterType() ExporterType {
	return ExporterType(c.Exporter)
}
