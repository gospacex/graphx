package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gospacex/graphx/observability"
)

func config() *observability.Config {
	return &observability.Config{
		Enabled:  true,
		Service:  "graphx-example",
		Endpoint: "localhost:4318",
		Protocol: "http",
		Exporter: "jaeger",
	}
}

func main() {
	ctx := context.Background()
	cfg := config()
	if err := observability.SetupTracing(ctx, cfg); err != nil {
		log.Fatalf("SetupTracing: %v", err)
	}
	defer observability.ShutdownTracerProvider(ctx)

	tr := observability.NewTracer()
	cctx, span := observability.StartQuerySpan(ctx, tr, "neo4j", "cypher", "MATCH (n) RETURN n", nil)
	observability.EndSpan(span, nil)

	fmt.Println("trace exported to Jaeger via OTLP HTTP", cctx)
}
