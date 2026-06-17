# graphx

**English** | [中文](./README.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/gospacex/graphx.svg)](https://pkg.go.dev/github.com/gospacex/graphx)
[![Go Report Card](https://goreportcard.com/badge/github.com/gospacex/graphx)](https://goreportcard.com/report/github.com/gospacex/graphx)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

> Enterprise-grade Go SDK for **7 graph databases** — Neo4j, Memgraph, Apache AGE, Dgraph, JanusGraph, NebulaGraph, TigerGraph.

---

## 📖 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Supported Databases](#-supported-databases)
- [Quick Start](#-quick-start)
- [Observability](#-observability)
- [Configuration Reference](#-configuration-reference)
- [Project Structure](#-project-structure)
- [Development Guide](#-development-guide)
- [Changelog](#-changelog)
- [License](#-license)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| **Multi-Database** | 7 mainstream graph databases, unified SDK interface |
| **Native Drivers** | Each backend returns its native driver handle — zero performance overhead |
| **Distributed Tracing** | OpenTelemetry integration with Jaeger / Kafka / Redis Stream exporters |
| **Configuration** | YAML config files + 28 quick functions, up and running in 5 minutes |
| **Production-Ready** | Connection pooling, retry mechanisms, timeout control, health checks |
| **No Unified Interface** | Preserves each database's native capabilities — no abstraction layer overhead |

---

## 🏗️ Architecture

```
graphx/
├── agex/              Apache AGE
├── dgraphx/           Dgraph
├── hubx/              Unified Configuration Hub
├── janusx/            JanusGraph
├── memgraphx/         Memgraph
├── neo4jx/            Neo4j
├── nebulax/           NebulaGraph
├── tigergraphx/       TigerGraph
├── tugraphx/          TuGraph
├── observability/     OpenTelemetry Observability
├── quick/             28 Quick Functions
└── internal/          Internal Utilities
```

**Two responsibility lines, fully decoupled:**

- **Backend subpackages** (`neo4jx/`, `memgraphx/`, etc.) — connection lifecycle management + native driver access
- **Observability** (`observability/`) — OTel distributed tracing with Jaeger / Kafka / Redis Stream exporters

> ⚠️ **Design Philosophy**: No unified `GraphDB` interface. Each backend returns its native driver handle, avoiding performance overhead and capability limitations from abstraction layers.

---

## 🗄️ Supported Databases

| Package | Database | Protocol | Driver |
|---------|----------|----------|--------|
| `neo4jx` | [Neo4j](https://neo4j.com/) | Bolt | `neo4j/neo4j-go-driver/v5` |
| `memgraphx` | [Memgraph](https://memgraph.com/) | Bolt | `neo4j/neo4j-go-driver/v5` |
| `agex` | [Apache AGE](https://age.apache.org/) | PostgreSQL | `lib/pq` + `database/sql` |
| `dgraphx` | [Dgraph](https://dgraph.io/) | HTTP/gRPC | Native `net/http` |
| `janusx` | [JanusGraph](https://janusgraph.org/) | Gremlin WebSocket | `apache/tinkerpop/gremlin-go/v3` |
| `nebulax` | [NebulaGraph](https://nebula-graph.io/) | Native | `vesoft-inc/nebula-go/v3` |
| `tigergraphx` | [TigerGraph](https://www.tigergraph.com/) | REST | Native `net/http` |
| `tugraphx` | [TuGraph](https://www.tugraphdb.com/) | RPC | Native RPC client |

---

## 🚀 Quick Start

### Method 1: Direct Initialization (Recommended)

```go
package main

import (
    "context"
    "github.com/gospacex/graphx/neo4jx"
)

func main() {
    ctx := context.Background()
    
    db, err := neo4jx.New(ctx, neo4jx.Config{
        Address:  "localhost:7687",
        Username: "neo4j",
        Password: "password",
        Database: "neo4j",
    })
    if err != nil {
        panic(err)
    }
    defer db.Close(ctx)

    // Access the native driver for full native API support
    driver := db.Driver() // neo4j.DriverWithContext
}
```

### Method 2: 28 Quick Functions

```go
package main

import (
    "github.com/gospacex/graphx/quick"
)

func main() {
    // One-liner with automatic YAML config parsing
    db, err := quick.N4PS("graphx.yaml")  // Neo4j + Memgraph
    if err != nil {
        panic(err)
    }
    defer db.Close(context.Background())
}
```

### Configuration File Example

```yaml
# graphx.yaml
neo4j:
  address: "localhost:7687"
  username: "neo4j"
  password: "password"
  database: "neo4j"
  max_connection_pool_size: 100
  connection_acquisition_timeout: 30s

memgraph:
  address: "localhost:7687"
  username: "memgraph"
  password: ""

agex:
  dsn: "postgres://postgres:password@localhost:5432/age?sslmode=disable"

dgraph:
  address: "localhost:8080"
  alpha_address: "localhost:9080"

janus:
  remote_hosts: "localhost:8182"
  read_timeout: 30s
  write_timeout: 30s

nebula:
  address: "localhost:9600"
  meta_address: "localhost:9559"
  username: "root"
  password: "nebula"

tigergraph:
  address: "localhost:9000"
  username: "digman"
  password: "digman"
  conn_timeout: 10s
  read_timeout: 30s
```

---

## 🔍 Observability

### Tracing Setup

```go
package main

import (
    "context"
    "github.com/gospacex/graphx/observability"
)

func main() {
    ctx := context.Background()
    
    err := observability.SetupTracing(ctx, &observability.Config{
        Enabled:  true,
        Service:  "my-graph-service",
        Exporter: "jaeger",           // jaeger | kafkatopic | redisstream
        Endpoint: "localhost:4318",   // OTLP HTTP endpoint
        // Or:
        // Exporter: "kafkatopic",
        // Endpoint: "localhost:9092",
        // KafkaTopic: "traces",
        // Or:
        // Exporter: "redisstream",
        // Endpoint: "localhost:6379",
        // RedisStreamKey: "traces",
    })
    if err != nil {
        panic(err)
    }
    defer observability.ShutdownTracerProvider(ctx)
}
```

### Trace Exporters

| Exporter | Protocol | Use Case |
|----------|----------|----------|
| `jaeger` | OTLP HTTP/gRPC | Dev/test environments, Jaeger UI |
| `kafkatopic` | Kafka Producer | Production, async decoupling |
| `redisstream` | Redis Stream | Lightweight, Redis ecosystem |

### Trace Output

Every database operation automatically generates spans:

```
📦 my-graph-service
 └─ 📊 neo4jx.Query
     ├─ 🔗 Session Acquire (12ms)
     ├─ 📝 Cypher: MATCH (n:User) RETURN n (45ms)
     └─ 🔗 Session Release (2ms)
```

---

## ⚙️ Configuration Reference

### Config Structure

```go
// Common config fields across all backends
type Config struct {
    Address                  string        `yaml:"address"`
    Username                 string        `yaml:"username"`
    Password                 string        `yaml:"password"`
    Database                 string        `yaml:"database"`
    MaxConnectionPoolSize    int           `yaml:"max_connection_pool_size"`
    ConnectionAcquisitionTimeout time.Duration `yaml:"connection_acquisition_timeout"`
    ConnTimeout              time.Duration `yaml:"conn_timeout"`
    ReadTimeout              time.Duration `yaml:"read_timeout"`
    WriteTimeout             time.Duration `yaml:"write_timeout"`
    SSL                      bool          `yaml:"ssl"`
    SSLInsecureSkipVerify    bool          `yaml:"ssl_insecure_skip_verify"`
}
```

### Backend-Specific Config

| Backend | Specific Fields |
|---------|-----------------|
| `neo4jx` | `DriverConfig` — full `neo4j.DriverWithContext` configuration |
| `agex` | `DSN` — PostgreSQL connection string |
| `dgraphx` | `AlphaAddress` — Dgraph Alpha address |
| `janusx` | `RemoteHosts`, `ReadTimeout`, `WriteTimeout` |
| `nebulax` | `MetaAddress` — Nebula Meta Server address |
| `tigergraphx` | `ConnTimeout`, `ReadTimeout` |

---

## 📁 Project Structure

```
graphx/
├── agex/                    Apache AGE backend
│   ├── db.go               Database connection wrapper
│   └── config.go           Config parsing
├── dgraphx/                 Dgraph backend
│   ├── client.go           HTTP/gRPC client
│   └── config.go
├── hubx/                    Unified Configuration Hub
│   ├── hub.go              Config aggregation & distribution
│   └── config.go
├── janusx/                  JanusGraph backend
│   ├── client.go           Gremlin WebSocket client
│   └── config.go
├── memgraphx/               Memgraph backend
│   ├── db.go               Bolt protocol wrapper (reuses neo4j driver)
│   └── config.go
├── neo4jx/                  Neo4j backend
│   ├── db.go               Bolt protocol wrapper
│   ├── config.go           Config parsing
│   └── driver.go           Native driver exposure
├── nebulax/                 NebulaGraph backend
│   ├── client.go           Native protocol wrapper
│   └── config.go
├── tigergraphx/             TigerGraph backend
│   ├── client.go           REST API wrapper
│   ├── jwt.go              JWT auth management
│   └── config.go
├── tugraphx/                TuGraph backend
│   ├── client.go           RPC client wrapper
│   └── config.go
├── observability/           Observability module
│   ├── tracing.go          OTel tracing initialization
│   ├── jaeger.go           Jaeger exporter
│   ├── kafka.go            Kafka exporter
│   └── redis.go            Redis Stream exporter
├── quick/                   Quick functions (28 total)
│   ├── quick.go            N4PS, N4PC, MGPS, etc.
│   └── config.go           Config loading
├── internal/                Internal utilities
│   └── util/               Shared utility functions
├── example/                 Usage examples
├── config.go                Root config structure
├── runtime.go               Runtime state management
├── Makefile                 Build scripts
└── go.mod                   Go module definition
```

---

## 🛠️ Development Guide

### Requirements

- **Go** ≥ 1.26.2
- **Git** ≥ 2.0

### Local Development

```bash
# Clone the repository
git clone https://github.com/gospacex/graphx.git
cd graphx

# Install dependencies
go mod download

# Run tests
go test ./... -v

# Format code
go fmt ./...

# Static analysis
go vet ./...

# Build
go build ./...
```

### Makefile Commands

```bash
make test       # Run all tests
make lint       # Lint code
make fmt        # Format code
make clean      # Clean build artifacts
```

### Adding a New Backend

1. Create a new subpackage directory at the root (e.g., `newdbx/`)
2. Implement the `DB` interface (`Connect`, `Close`, `Driver`)
3. Add the `Config` struct to `config.go`
4. Add a quick function to `quick/quick.go`
5. Update the Supported Databases table in this README

---

## 📋 Changelog

| Version | Date | Description |
|---------|------|-------------|
| v1.0.0 | 2026-06-18 | Initial release — 7 database backends + observability |

[CHANGELOG.md](./CHANGELOG.md)

---

## 📄 License

This project is licensed under the **MIT License**. See [LICENSE](./LICENSE) for details.

---

## 🤝 Contributing

We welcome Issues and Pull Requests!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📧 Contact

- **Repository**: https://github.com/gospacex/graphx
- **Issues**: https://github.com/gospacex/graphx/issues

---

*Made with ❤️ by the graphx team*
