package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gospacex/graphx/internal/mqxbinding"
	"github.com/gospacex/graphx/observability"
	"github.com/gospacex/graphx/observability/exporter/kafkatopic"
)

func main() {
	ctx := context.Background()

	mqCfg := mqxbinding.MQConfig{
		Driver: "kafka",
		Mode:   "cluster",
		Addrs:  []string{"localhost:9092"},
	}
	prod, err := mqxbinding.Kafka(mqCfg)
	if err != nil {
		log.Fatalf("mqxbinding.Kafka: %v", err)
	}

	cfg := &observability.Config{
		Enabled:    true,
		Service:    "graphx-example",
		Exporter:   "kafkatopic",
		KafkaTopic: "otel-traces",
	}
	if err := observability.SetupTracing(ctx, cfg); err != nil {
		log.Fatalf("SetupTracing: %v", err)
	}
	defer observability.ShutdownTracerProvider(ctx)

	_ = prod
	_ = kafkatopic.New
	fmt.Println("trace exported to Kafka topic otel-traces")
}
