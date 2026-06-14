// Package agex provides an Apache AGE backend via PostgreSQL lib/pq driver.
package agex

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/observability"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// AGE wraps a *sql.DB and exposes lifecycle + accessor methods for Apache AGE.
type AGE struct {
	db     *sql.DB
	cfg    graphx.Config
	graph  string
	tracer trace.Tracer
	mu     sync.RWMutex
}

// New creates an AGE driver, verifies connectivity, and returns a ready handle.
func New(ctx context.Context, cfg graphx.Config) (*AGE, error) {
	if ctx == nil {
		slog.Warn("graphx/agex: nil ctx, using context.Background()")
		ctx = context.Background()
	}
	tracer := observability.TracerForBackend("agex", cfg.TracerName)
	ctx, span := tracer.Start(ctx, "graphx.agex.new",
		trace.WithAttributes(
			attribute.String("db.system", "age"),
			attribute.String("server.address", cfg.Address),
			attribute.String("db.operation", "new"),
		))
	defer span.End()

	acfg, err := Build(cfg)
	if err != nil {
		recordSpanError(span, err)
		return nil, err
	}

	db, err := sql.Open("postgres", acfg.DSN)
	if err != nil {
		recordSpanError(span, err)
		return nil, fmt.Errorf("graphx/agex: %w", err)
	}

	if cfg.PoolSize > 0 {
		db.SetMaxOpenConns(cfg.PoolSize)
		db.SetMaxIdleConns(cfg.PoolSize / 2)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		recordSpanError(span, err)
		return nil, fmt.Errorf("graphx/agex: ping: %w", err)
	}

	return &AGE{db: db, cfg: cfg, graph: acfg.Graph, tracer: tracer}, nil
}

// Close implements graphx.Closer.
func (g *AGE) Close(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.db == nil {
		return nil
	}
	_, span := g.tracer.Start(ctx, "graphx.agex.close",
		trace.WithAttributes(
			attribute.String("db.system", "age"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "close"),
		))
	defer span.End()
	if err := g.db.Close(); err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("graphx/agex: close: %w", err)
	}
	g.db = nil
	return nil
}

// Ping implements graphx.Pinger.
func (g *AGE) Ping(ctx context.Context) error {
	g.mu.RLock()
	db := g.db
	g.mu.RUnlock()
	if db == nil {
		return fmt.Errorf("graphx/agex: ping: closed")
	}
	_, span := g.tracer.Start(ctx, "graphx.agex.ping",
		trace.WithAttributes(
			attribute.String("db.system", "age"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "ping"),
		))
	defer span.End()
	if err := db.PingContext(ctx); err != nil {
		slog.Warn("graphx/agex: ping failed", "address", g.cfg.Address, "error", err)
		recordSpanError(span, err)
		return fmt.Errorf("graphx/agex: ping: %w", err)
	}
	slog.Debug("graphx/agex: ping ok", "address", g.cfg.Address)
	return nil
}

func recordSpanError(span trace.Span, err error) {
	if err == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// DB returns the underlying *sql.DB for direct use.
func (g *AGE) DB() *sql.DB {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.db
}

// Graph returns the AGE graph name configured for this instance.
func (g *AGE) Graph() string { return g.graph }

// Cfg returns the config used to create this instance.
func (g *AGE) Cfg() graphx.Config { return g.cfg }

// --- singleton ---

var (
	defaultInst = graphx.NewSingleton[AGE]("agex")
	namedMu     sync.Mutex
	namedInst   = map[string]*AGE{}
)

// Get returns the shared singleton AGE instance, creating it on first call.
// When cfg.Name is non-empty the instance is keyed by name.
func Get(ctx context.Context, cfg graphx.Config) (*AGE, error) {
	if cfg.Name != "" {
		return getNamed(ctx, cfg)
	}
	return defaultInst.Get(func() (*AGE, error) { return New(ctx, cfg) })
}

func getNamed(ctx context.Context, cfg graphx.Config) (*AGE, error) {
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
	namedInst = map[string]*AGE{}
	namedMu.Unlock()
}
