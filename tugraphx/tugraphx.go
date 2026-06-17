package tugraphx

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

// TuGraph wraps a neo4j.DriverWithContext (TuGraph Bolt/Cypher) and exposes lifecycle + accessor methods.
type TuGraph struct {
	driver neo4j.DriverWithContext
	cfg    graphx.Config
	tracer trace.Tracer
	mu     sync.RWMutex
}

// New creates a TuGraph Bolt driver, verifies connectivity, and returns a ready handle.
func New(ctx context.Context, cfg graphx.Config) (*TuGraph, error) {
	if ctx == nil {
		slog.Warn("graphx/tugraphx: nil ctx, using context.Background()")
		ctx = context.Background()
	}
	tracer := observability.TracerForBackend("tugraphx", cfg.TracerName)
	ctx, span := tracer.Start(ctx, "graphx.tugraphx.new",
		trace.WithAttributes(
			attribute.String("db.system", "tugraph"),
			attribute.String("server.address", cfg.Address),
			attribute.String("db.operation", "new"),
		))
	defer span.End()

	tcfg, err := Build(cfg)
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

	driver, err := neo4j.NewDriverWithContext(tcfg.BoltURI, auth, func(c *neo4j.Config) {
		if cfg.PoolSize > 0 {
			c.MaxConnectionPoolSize = cfg.PoolSize
		}
	})
	if err != nil {
		recordSpanError(span, err)
		return nil, fmt.Errorf("graphx/tugraphx: %w", err)
	}

	if err := driver.VerifyConnectivity(ctx); err != nil {
		_ = driver.Close(ctx)
		recordSpanError(span, err)
		return nil, fmt.Errorf("graphx/tugraphx: %w", err)
	}

	return &TuGraph{driver: driver, cfg: cfg, tracer: tracer}, nil
}

// Close implements graphx.Closer. Repeated calls are safe.
func (g *TuGraph) Close(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.driver == nil {
		return nil
	}
	_, span := g.tracer.Start(ctx, "graphx.tugraphx.close",
		trace.WithAttributes(
			attribute.String("db.system", "tugraph"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "close"),
		))
	defer span.End()
	if err := g.driver.Close(ctx); err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("graphx/tugraphx: close: %w", err)
	}
	g.driver = nil
	return nil
}

// Ping implements graphx.Pinger.
func (g *TuGraph) Ping(ctx context.Context) error {
	g.mu.RLock()
	drv := g.driver
	g.mu.RUnlock()
	if drv == nil {
		return fmt.Errorf("graphx/tugraphx: ping: closed")
	}
	_, span := g.tracer.Start(ctx, "graphx.tugraphx.ping",
		trace.WithAttributes(
			attribute.String("db.system", "tugraph"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "ping"),
		))
	defer span.End()
	if err := drv.VerifyConnectivity(ctx); err != nil {
		slog.Warn("graphx/tugraphx: ping failed", "address", g.cfg.Address, "error", err)
		recordSpanError(span, err)
		return fmt.Errorf("graphx/tugraphx: ping: %w", err)
	}
	slog.Debug("graphx/tugraphx: ping ok", "address", g.cfg.Address)
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
func (g *TuGraph) Driver() neo4j.DriverWithContext {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.driver
}

// Cfg returns the config used to create this instance.
func (g *TuGraph) Cfg() graphx.Config { return g.cfg }

// --- singleton ---

var (
	defaultInst = graphx.NewSingleton[TuGraph]("tugraphx")
	namedMu     sync.Mutex
	namedInst   = map[string]*TuGraph{}
)

// Get returns the shared singleton TuGraph instance, creating it on first call.
// When cfg.Name is non-empty the instance is keyed by name.
func Get(ctx context.Context, cfg graphx.Config) (*TuGraph, error) {
	if cfg.Name != "" {
		return getNamed(ctx, cfg)
	}
	return defaultInst.Get(func() (*TuGraph, error) { return New(ctx, cfg) })
}

func getNamed(ctx context.Context, cfg graphx.Config) (*TuGraph, error) {
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
	namedInst = map[string]*TuGraph{}
	namedMu.Unlock()
}
