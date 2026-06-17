// Package memgraphx provides a Memgraph backend via Bolt/Cypher using neo4j-go-driver.
package memgraphx

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/observability"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Memgraph wraps a neo4j.DriverWithContext and exposes lifecycle + accessor methods.
type Memgraph struct {
	driver neo4j.DriverWithContext
	cfg    graphx.Config
	tracer trace.Tracer
	mu     sync.RWMutex
}

// New creates a Memgraph Bolt driver, verifies connectivity, and returns a ready handle.
func New(ctx context.Context, cfg graphx.Config) (*Memgraph, error) {
	if ctx == nil {
		slog.Warn("graphx/memgraphx: nil ctx, using context.Background()")
		ctx = context.Background()
	}
	tracer := observability.TracerForBackend("memgraphx", cfg.TracerName)
	ctx, span := tracer.Start(ctx, "graphx.memgraphx.new",
		trace.WithAttributes(
			attribute.String("db.system", "memgraph"),
			attribute.String("server.address", cfg.Address),
			attribute.String("db.operation", "new"),
		))
	defer span.End()

	mcfg, err := Build(cfg)
	if err != nil {
		recordSpanError(span, err)
		return nil, err
	}

	var auth neo4j.AuthToken
	if cfg.Username != "" {
		auth = neo4j.BasicAuth(cfg.Username, cfg.Password, "")
	} else {
		auth = neo4j.NoAuth()
	}

	driver, err := neo4j.NewDriverWithContext(mcfg.BoltURI, auth, func(c *neo4j.Config) {
		if cfg.PoolSize > 0 {
			c.MaxConnectionPoolSize = cfg.PoolSize
		}
	})
	if err != nil {
		recordSpanError(span, err)
		return nil, fmt.Errorf("graphx/memgraphx: %w", err)
	}

	if err := driver.VerifyConnectivity(ctx); err != nil {
		_ = driver.Close(ctx)
		recordSpanError(span, err)
		return nil, fmt.Errorf("graphx/memgraphx: %w", err)
	}

	return &Memgraph{driver: driver, cfg: cfg, tracer: tracer}, nil
}

// Close implements graphx.Closer. Repeated calls are safe.
func (g *Memgraph) Close(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.driver == nil {
		return nil
	}
	_, span := g.tracer.Start(ctx, "graphx.memgraphx.close",
		trace.WithAttributes(
			attribute.String("db.system", "memgraph"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "close"),
		))
	defer span.End()
	if err := g.driver.Close(ctx); err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("graphx/memgraphx: close: %w", err)
	}
	g.driver = nil
	return nil
}

// Ping implements graphx.Pinger by executing RETURN 1.
func (g *Memgraph) Ping(ctx context.Context) error {
	g.mu.RLock()
	drv := g.driver
	g.mu.RUnlock()
	if drv == nil {
		return fmt.Errorf("graphx/memgraphx: ping: closed")
	}
	_, span := g.tracer.Start(ctx, "graphx.memgraphx.ping",
		trace.WithAttributes(
			attribute.String("db.system", "memgraph"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "ping"),
		))
	defer span.End()
	session := drv.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.Run(ctx, "RETURN 1", nil)
	if err != nil {
		slog.Warn("graphx/memgraphx: ping failed", "address", g.cfg.Address, "error", err)
		recordSpanError(span, err)
		return fmt.Errorf("graphx/memgraphx: ping: %w", err)
	}
	slog.Debug("graphx/memgraphx: ping ok", "address", g.cfg.Address)
	return nil
}

func recordSpanError(span trace.Span, err error) {
	if err == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// Driver returns the underlying neo4j.DriverWithContext for direct use.
func (g *Memgraph) Driver() neo4j.DriverWithContext {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.driver
}

// Cfg returns the config used to create this instance.
func (g *Memgraph) Cfg() graphx.Config { return g.cfg }

// --- singleton ---

var (
	defaultInst = graphx.NewSingleton[Memgraph]("memgraphx")
	namedMu     sync.Mutex
	namedInst   = map[string]*Memgraph{}
)

// Get returns the shared singleton Memgraph instance, creating it on first call.
// When cfg.Name is non-empty the instance is keyed by name.
func Get(ctx context.Context, cfg graphx.Config) (*Memgraph, error) {
	if cfg.Name != "" {
		return getNamed(ctx, cfg)
	}
	return defaultInst.Get(func() (*Memgraph, error) { return New(ctx, cfg) })
}

func getNamed(ctx context.Context, cfg graphx.Config) (*Memgraph, error) {
	namedMu.Lock()
	defer namedMu.Unlock()
	if existing, ok := namedInst[cfg.Name]; ok {
		return existing, nil
	}
	inst, err := New(ctx, cfg)
	if err != nil {
		return nil, err
	}
	namedInst[cfg.Name] = inst
	return inst, nil
}

// Reset clears all cached instances (testing only).
func Reset() {
	defaultInst.Reset()
	namedMu.Lock()
	for _, inst := range namedInst {
		_ = inst.Close(context.Background())
	}
	namedInst = map[string]*Memgraph{}
	namedMu.Unlock()
}
