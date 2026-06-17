# graphx

[English](./README.en.md) | **中文**

[![Go Reference](https://pkg.go.dev/badge/github.com/gospacex/graphx.svg)](https://pkg.go.dev/github.com/gospacex/graphx)
[![Go Report Card](https://goreportcard.com/badge/github.com/gospacex/graphx)](https://goreportcard.com/report/github.com/gospacex/graphx)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

> 企业级 Go SDK，统一接入 **7 款图数据库**：Neo4j、Memgraph、Apache AGE、Dgraph、JanusGraph、NebulaGraph、TigerGraph。

---

## 📖 目录

- [特性](#-特性)
- [架构设计](#-架构设计)
- [支持的数据库](#-支持的数据库)
- [快速开始](#-快速开始)
- [可观测性](#-可观测性)
- [配置参考](#-配置参考)
- [项目结构](#-项目结构)
- [开发指南](#-开发指南)
- [版本历史](#-版本历史)
- [许可证](#-许可证)

---

## ✨ 特性

| 特性 | 说明 |
|------|------|
| **多数据库支持** | 7 款主流图数据库，统一 SDK 接口 |
| **原生驱动** | 各后端返回原生驱动句柄，零性能损耗 |
| **分布式追踪** | OpenTelemetry 集成，支持 Jaeger / Kafka / Redis Stream |
| **配置管理** | YAML 配置文件 + 28 个快捷函数，5 分钟上手 |
| **生产就绪** | 连接池管理、重试机制、超时控制、健康检查 |
| **无统一接口** | 保持各数据库原生能力，不做抽象层损耗 |

---

## 🏗️ 架构设计

```
graphx/
├── agex/              Apache AGE
├── dgraphx/           Dgraph
├── hubx/              统一配置中心
├── janusx/            JanusGraph
├── memgraphx/         Memgraph
├── neo4jx/            Neo4j
├── nebulax/           NebulaGraph
├── tigergraphx/       TigerGraph
├── tugraphx/          TuGraph
├── observability/     OpenTelemetry 可观测性
├── quick/             28 个快捷函数
└── internal/          内部工具
```

**两大职责线，完全解耦：**

- **后端子包**（`neo4jx/`、`memgraphx/` 等）—— 连接生命周期管理 + 原生驱动访问
- **可观测性**（`observability/`）—— OTel 链路追踪，支持 Jaeger / Kafka / Redis Stream 导出

> ⚠️ **设计哲学**：不提供统一的 `GraphDB` 接口。每个后端返回其原生驱动句柄，避免抽象层带来的性能损耗和能力限制。

---

## 🗄️ 支持的数据库

| 子包 | 数据库 | 协议 | 驱动 |
|------|--------|------|------|
| `neo4jx` | [Neo4j](https://neo4j.com/) | Bolt | `neo4j/neo4j-go-driver/v5` |
| `memgraphx` | [Memgraph](https://memgraph.com/) | Bolt | `neo4j/neo4j-go-driver/v5` |
| `agex` | [Apache AGE](https://age.apache.org/) | PostgreSQL | `lib/pq` + `database/sql` |
| `dgraphx` | [Dgraph](https://dgraph.io/) | HTTP/gRPC | 原生 `net/http` |
| `janusx` | [JanusGraph](https://janusgraph.org/) | Gremlin WebSocket | `apache/tinkerpop/gremlin-go/v3` |
| `nebulax` | [NebulaGraph](https://nebula-graph.io/) | Native | `vesoft-inc/nebula-go/v3` |
| `tigergraphx` | [TigerGraph](https://www.tigergraph.com/) | REST | 原生 `net/http` |
| `tugraphx` | [TuGraph](https://www.tugraphdb.com/) | RPC | 原生 RPC 客户端 |

---

## 🚀 快速开始

### 方式一：直接初始化（推荐）

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

    // 获取原生驱动，使用全部原生 API
    driver := db.Driver() // neo4j.DriverWithContext
}
```

### 方式二：28 个快捷函数

```go
package main

import (
    "github.com/gospacex/graphx/quick"
)

func main() {
    // 一行代码，自动解析 YAML 配置
    db, err := quick.N4PS("graphx.yaml")  // Neo4j + Memgraph
    if err != nil {
        panic(err)
    }
    defer db.Close(context.Background())
}
```

### 配置文件示例

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

## 🔍 可观测性

### 链路追踪 Setup

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
        // 或
        // Exporter: "kafkatopic",
        // Endpoint: "localhost:9092",
        // KafkaTopic: "traces",
        // 或
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

### 追踪导出器

| 导出器 | 协议 | 用途 |
|--------|------|------|
| `jaeger` | OTLP HTTP/gRPC | 开发/测试环境，Jaeger UI |
| `kafkatopic` | Kafka Producer | 生产环境，异步解耦 |
| `redisstream` | Redis Stream | 轻量级，Redis 生态 |

### 追踪效果

每个数据库操作自动生成 span：

```
📦 my-graph-service
 └─ 📊 neo4jx.Query
     ├─ 🔗 Session Acquire (12ms)
     ├─ 📝 Cypher: MATCH (n:User) RETURN n (45ms)
     └─ 🔗 Session Release (2ms)
```

---

## ⚙️ 配置参考

### Config 结构

```go
// 各后端通用配置字段
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

### 各后端特有配置

| 后端 | 特有配置 |
|------|----------|
| `neo4jx` | `DriverConfig` — 完整 neo4j.DriverWithContext 配置 |
| `agex` | `DSN` — PostgreSQL 连接字符串 |
| `dgraphx` | `AlphaAddress` — Dgraph Alpha 地址 |
| `janusx` | `RemoteHosts`, `ReadTimeout`, `WriteTimeout` |
| `nebulax` | `MetaAddress` — Nebula Meta Server 地址 |
| `tigergraphx` | `ConnTimeout`, `ReadTimeout` |

---

## 📁 项目结构

```
graphx/
├── agex/                    Apache AGE 后端
│   ├── db.go               数据库连接封装
│   └── config.go           配置解析
├── dgraphx/                 Dgraph 后端
│   ├── client.go           HTTP/gRPC 客户端
│   └── config.go
├── hubx/                    统一配置中心
│   ├── hub.go              配置聚合与分发
│   └── config.go
├── janusx/                  JanusGraph 后端
│   ├── client.go           Gremlin WebSocket 客户端
│   └── config.go
├── memgraphx/               Memgraph 后端
│   ├── db.go               Bolt 协议封装（复用 neo4j 驱动）
│   └── config.go
├── neo4jx/                  Neo4j 后端
│   ├── db.go               Bolt 协议封装
│   ├── config.go           配置解析
│   └── driver.go           原生驱动暴露
├── nebulax/                 NebulaGraph 后端
│   ├── client.go           Native 协议封装
│   └── config.go
├── tigergraphx/             TigerGraph 后端
│   ├── client.go           REST API 封装
│   ├── jwt.go              JWT 认证管理
│   └── config.go
├── tugraphx/                TuGraph 后端
│   ├── client.go           RPC 客户端封装
│   └── config.go
├── observability/           可观测性模块
│   ├── tracing.go          OTel 追踪初始化
│   ├── jaeger.go           Jaeger 导出器
│   ├── kafka.go            Kafka 导出器
│   └── redis.go            Redis Stream 导出器
├── quick/                   快捷函数（28 个）
│   ├── quick.go            N4PS, N4PC, MGPS 等
│   └── config.go           配置加载
├── internal/                内部工具
│   └── util/               共享工具函数
├── example/                 使用示例
├── config.go                根配置结构
├── runtime.go               运行时状态管理
├── Makefile                 构建脚本
└── go.mod                   Go 模块定义
```

---

## 🛠️ 开发指南

### 环境要求

- **Go** ≥ 1.26.2
- **Git** ≥ 2.0

### 本地开发

```bash
# 克隆仓库
git clone https://github.com/gospacex/graphx.git
cd graphx

# 安装依赖
go mod download

# 运行测试
go test ./... -v

# 格式化代码
go fmt ./...

# 静态分析
go vet ./...

# 构建
go build ./...
```

### Makefile 命令

```bash
make test       # 运行所有测试
make lint       # 代码检查
make fmt        # 代码格式化
make clean      # 清理构建产物
```

### 添加新后端

1. 在根目录创建子包目录（如 `newdbx/`）
2. 实现 `DB` 接口（`Connect`、`Close`、`Driver`）
3. 添加 `Config` 结构到 `config.go`
4. 在 `quick/quick.go` 添加快捷函数
5. 在本文档更新支持的数据库表格

---

## 📋 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| v1.0.0 | 2026-06-18 | 初始发布，7 款数据库后端 + 可观测性 |

[CHANGELOG.md](./CHANGELOG.md)

---

## 📄 许可证

本项目采用 **MIT 许可证**。详见 [LICENSE](./LICENSE)。

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📧 联系方式

- **项目地址**: https://github.com/gospacex/graphx
- **问题反馈**: https://github.com/gospacex/graphx/issues

---

*Made with ❤️ by the graphx team*
