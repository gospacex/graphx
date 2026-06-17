// Package nebulax provides a NebulaGraph backend via the native nebula-go SessionPool.
package nebulax

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/observability"
	ng "github.com/vesoft-inc/nebula-go/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Nebula wraps a *ng.SessionPool and exposes lifecycle + accessor methods.
type Nebula struct {
	pool   *ng.SessionPool
	cfg    graphx.Config
	space  string
	tracer trace.Tracer
	mu     sync.RWMutex
}

// New creates a NebulaGraph session pool and returns a ready handle.
func New(ctx context.Context, cfg graphx.Config) (*Nebula, error) {
	if ctx == nil {
		slog.Warn("graphx/nebulax: nil ctx, using context.Background()")
		ctx = context.Background()
	}
	tracer := observability.TracerForBackend("nebulax", cfg.TracerName)
	ctx, span := tracer.Start(ctx, "graphx.nebulax.new",
		trace.WithAttributes(
			attribute.String("db.system", "nebula"),
			attribute.String("server.address", cfg.Address),
			attribute.String("db.operation", "new"),
		))
	defer span.End()

	ncfg, err := Build(cfg)
	if err != nil {
		recordSpanError(span, err)
		return nil, err
	}

	hosts := make([]ng.HostAddress, len(ncfg.Hosts))
	for i, h := range ncfg.Hosts {
		host, port := splitHostPort(h)
		hosts[i] = ng.HostAddress{Host: host, Port: port}
	}

	opts := []ng.SessionPoolConfOption{
		ng.WithMinSize(cfg.PoolSize / 2),
		ng.WithMaxSize(cfg.PoolSize),
	}
	if cfg.PoolSize <= 0 {
		opts = nil
	}
	if cfg.TLS {
		opts = append(opts, ng.WithSSLConfig(cfg.TLSConfig))
	}

	poolConf, err := ng.NewSessionPoolConf(
		cfg.Username, cfg.Password,
		hosts, ncfg.Space, opts...,
	)
	if err != nil {
		recordSpanError(span, err)
		return nil, fmt.Errorf("graphx/nebulax: %w", err)
	}

	pool, err := ng.NewSessionPool(*poolConf, ng.DefaultLogger{})
	if err != nil {
		recordSpanError(span, err)
		return nil, fmt.Errorf("graphx/nebulax: pool: %w", err)
	}

	return &Nebula{pool: pool, cfg: cfg, space: ncfg.Space, tracer: tracer}, nil
}

// Close implements graphx.Closer.
func (g *Nebula) Close(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.pool == nil {
		return nil
	}
	_, span := g.tracer.Start(ctx, "graphx.nebulax.close",
		trace.WithAttributes(
			attribute.String("db.system", "nebula"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "close"),
		))
	defer span.End()
	g.pool.Close()
	g.pool = nil
	return nil
}

// Ping implements graphx.Pinger by executing SHOW SPACES.
func (g *Nebula) Ping(ctx context.Context) error {
	g.mu.RLock()
	pool := g.pool
	g.mu.RUnlock()
	if pool == nil {
		return fmt.Errorf("graphx/nebulax: ping: closed")
	}
	_, span := g.tracer.Start(ctx, "graphx.nebulax.ping",
		trace.WithAttributes(
			attribute.String("db.system", "nebula"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "ping"),
		))
	defer span.End()
	result, err := pool.ExecuteAndCheck("SHOW SPACES")
	if err != nil {
		slog.Warn("graphx/nebulax: ping failed", "address", g.cfg.Address, "error", err)
		recordSpanError(span, err)
		return fmt.Errorf("graphx/nebulax: ping: %w", err)
	}
	if result.IsSucceed() {
		slog.Debug("graphx/nebulax: ping ok", "address", g.cfg.Address)
		return nil
	}
	err = fmt.Errorf("graphx/nebulax: ping: %s", result.GetErrorMsg())
	recordSpanError(span, err)
	return err
}

func recordSpanError(span trace.Span, err error) {
	if err == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// Pool returns the underlying *ng.SessionPool for direct use.
func (g *Nebula) Pool() *ng.SessionPool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.pool
}

// Space returns the NebulaGraph space name configured for this instance.
func (g *Nebula) Space() string { return g.space }

// Cfg returns the config used to create this instance.
func (g *Nebula) Cfg() graphx.Config { return g.cfg }

// --- singleton ---

var (
	defaultInst = graphx.NewSingleton[Nebula]("nebulax")
	namedMu     sync.Mutex
	namedInst   = map[string]*Nebula{}
)

// Get returns the shared singleton Nebula instance, creating it on first call.
// When cfg.Name is non-empty the instance is keyed by name.
func Get(ctx context.Context, cfg graphx.Config) (*Nebula, error) {
	if cfg.Name != "" {
		return getNamed(ctx, cfg)
	}
	return defaultInst.Get(func() (*Nebula, error) { return New(ctx, cfg) })
}

func getNamed(ctx context.Context, cfg graphx.Config) (*Nebula, error) {
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
	namedInst = map[string]*Nebula{}
	namedMu.Unlock()
}

// splitHostPort splits "host:port" into separate strings and int.
func splitHostPort(addr string) (string, int) {
	port := 9669
	host := addr
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			host = addr[:i]
			fmt.Sscanf(addr[i+1:], "%d", &port)
			break
		}
	}
	return host, port
}
