package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gospacex/graphx/internal/mqxbinding"
	"github.com/gospacex/graphx/observability"
	"github.com/gospacex/graphx/observability/exporter/redisstream"
)

func main() {
	ctx := context.Background()

	mqCfg := mqxbinding.MQConfig{
		Driver: "redis",
		Mode:   "single",
		Addrs:  []string{"localhost:6379"},
	}
	cli, err := mqxbinding.Redis(mqCfg)
	if err != nil {
		log.Fatalf("mqxbinding.Redis: %v", err)
	}

	cfg := &observability.Config{
		Enabled:     true,
		Service:     "graphx-example",
		Exporter:    "redisstream",
		RedisStream: "otel-traces",
	}
	if err := observability.SetupTracing(ctx, cfg); err != nil {
		log.Fatalf("SetupTracing: %v", err)
	}
	defer observability.ShutdownTracerProvider(ctx)

	_ = cli
	_ = redisstream.New
	fmt.Println("trace exported to Redis stream otel-traces")
}
