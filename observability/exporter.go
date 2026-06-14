package observability

import (
	"context"
	"fmt"
	"log/slog"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/gospacex/graphx/internal/mqxbinding"
	"github.com/gospacex/graphx/observability/exporter/jaeger"
	"github.com/gospacex/graphx/observability/exporter/kafkatopic"
	"github.com/gospacex/graphx/observability/exporter/redisstream"
)

func newSpanExporter(ctx context.Context, cfg *Config) (sdktrace.SpanExporter, error) {
	switch cfg.GetExporterType() {
	case ExporterJaeger:
		return jaeger.New(ctx, cfg.Endpoint, cfg.Protocol)
	case ExporterKafkaTopic:
		mqCfg := mqxbinding.MQConfig{
			Driver: "kafka",
			Mode:   "cluster",
			Addrs:  cfg.KafkaBrokers,
			Kafka:  mqxbinding.KafkaConfig{SecurityProtocol: cfg.KafkaSASL.Mechanism},
		}
		prod, err := mqxbinding.Kafka(mqCfg)
		if err != nil {
			return nil, fmt.Errorf("observability: kafka producer: %w", err)
		}
		slog.Debug("graphx/observability: using kafka exporter", "topic", cfg.KafkaTopic)
		return kafkatopic.New(prod, cfg.KafkaTopic), nil
	case ExporterRedisStream:
		mqCfg := mqxbinding.MQConfig{
			Driver: "redis",
			Mode:   "single",
			Addrs:  cfg.RedisAddrs,
			Redis:  mqxbinding.RedisConfig{DB: 0},
		}
		cli, err := mqxbinding.Redis(mqCfg)
		if err != nil {
			return nil, fmt.Errorf("observability: redis client: %w", err)
		}
		slog.Debug("graphx/observability: using redis exporter", "stream", cfg.RedisStream)
		return redisstream.New(cli, cfg.RedisStream), nil
	default:
		return nil, fmt.Errorf("observability: unknown exporter %q", cfg.Exporter)
	}
}

func newSampler(cfg *Config) sdktrace.Sampler {
	switch cfg.SamplerType {
	case "always_on":
		return sdktrace.AlwaysSample()
	case "always_off":
		return sdktrace.NeverSample()
	case "traceidratio":
		return sdktrace.TraceIDRatioBased(cfg.SamplerRatio)
	case "parentbased_always_on":
		return sdktrace.ParentBased(sdktrace.AlwaysSample())
	default:
		return sdktrace.AlwaysSample()
	}
}

func newResource(serviceName string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		attribute.String("service.name", serviceName),
	)
}
