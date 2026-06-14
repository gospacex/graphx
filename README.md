# graphx

Production-grade Go SDK for 7 graph databases — Neo4j, Memgraph, Apache AGE, Dgraph, JanusGraph, Nebula, TigerGraph.

## Architecture

Two responsibility lines, fully decoupled:

- **Backend subpackages** (`neo4jx/`, `memgraphx/`, ...) — connection lifecycle + native driver access
- **Observability** (`observability/`) — OTel tracing with Jaeger / Kafka / Redis stream exporters

No unified `GraphDB` interface. Each backend returns its native driver handle.

## Backends

| Package | Database | Driver |
|---------|----------|--------|
| `neo4jx` | Neo4j | `neo4j-go-driver/v5` (Bolt) |
| `memgraphx` | Memgraph | `neo4j-go-driver/v5` (Bolt) |
| `agex` | Apache AGE | `database/sql` + `lib/pq` |
| `dgraphx` | Dgraph | `net/http` REST |
| `janusx` | JanusGraph | `gremlin-go/v3` (WebSocket) |
| `nebulax` | NebulaGraph | `nebula-go/v3` |
| `tigergraphx` | TigerGraph | `net/http` REST |

## Quick Start

```go
import "github.com/gospacex/graphx/neo4jx"

db, _ := neo4jx.New(ctx, graphx.Config{
    Address: "localhost:7687",
    Username: "neo4j",
    Password: "password",
})
defer db.Close(ctx)
raw := db.Driver() // neo4j.DriverWithContext
```

Or with 28 quick functions:

```go
import "github.com/gospacex/graphx/quick"

db, _ := quick.N4PS("graphx.yaml")
```

## Tracing

```go
import "github.com/gospacex/graphx/observability"

observability.SetupTracing(ctx, &observability.Config{
    Enabled:  true,
    Service:  "my-service",
    Exporter: "jaeger",
    Endpoint: "localhost:4318",
})
defer observability.ShutdownTracerProvider(ctx)
```

Exporters: `jaeger` (OTLP HTTP/gRPC), `kafkatopic`, `redisstream`.

## License

MIT
