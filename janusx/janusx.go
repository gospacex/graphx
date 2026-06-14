// Package janusx provides a JanusGraph backend via the Gremlin WebSocket protocol.
package janusx

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/observability"
	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// JanusGraph wraps a *gremlingo.Client and exposes lifecycle + accessor methods.
type JanusGraph struct {
	client *gremlingo.Client
	cfg    graphx.Config
	tracer trace.Tracer
	mu     sync.RWMutex
}

// New creates a JanusGraph Gremlin client and returns a ready handle.
func New(ctx context.Context, cfg graphx.Config) (*JanusGraph, error) {
	if ctx == nil {
		slog.Warn("graphx/janusx: nil ctx, using context.Background()")
		ctx = context.Background()
	}
	tracer := observability.TracerForBackend("janusx", cfg.TracerName)
	ctx, span := tracer.Start(ctx, "graphx.janusx.new",
		trace.WithAttributes(
			attribute.String("db.system", "janusgraph"),
			attribute.String("server.address", cfg.Address),
			attribute.String("db.operation", "new"),
		))
	defer span.End()

	jcfg, err := Build(cfg)
	if err != nil {
		recordSpanError(span, err)
		return nil, err
	}

	var tlsCfg *tls.Config
	if cfg.TLS {
		tlsCfg = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	var auth gremlingo.AuthInfoProvider
	if cfg.Username != "" {
		auth = &gremlingo.AuthInfo{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	client, err := gremlingo.NewClient(jcfg.URL, func(settings *gremlingo.ClientSettings) {
		settings.TraversalSource = "g"
		settings.KeepAliveInterval = 30 * time.Second
		settings.TlsConfig = tlsCfg
		settings.AuthInfo = auth
	})
	if err != nil {
		recordSpanError(span, err)
		return nil, fmt.Errorf("graphx/janusx: %w", err)
	}

	return &JanusGraph{client: client, cfg: cfg, tracer: tracer}, nil
}

// Close implements graphx.Closer. Repeated calls are safe.
func (g *JanusGraph) Close(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.client == nil {
		return nil
	}
	_, span := g.tracer.Start(ctx, "graphx.janusx.close",
		trace.WithAttributes(
			attribute.String("db.system", "janusgraph"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "close"),
		))
	defer span.End()
	g.client.Close()
	g.client = nil
	return nil
}

// Ping implements graphx.Pinger by submitting g.V().limit(1).
func (g *JanusGraph) Ping(ctx context.Context) error {
	g.mu.RLock()
	cl := g.client
	g.mu.RUnlock()
	if cl == nil {
		return fmt.Errorf("graphx/janusx: ping: closed")
	}
	_, span := g.tracer.Start(ctx, "graphx.janusx.ping",
		trace.WithAttributes(
			attribute.String("db.system", "janusgraph"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "ping"),
		))
	defer span.End()
	result, err := cl.Submit("g.V().limit(1)")
	if err != nil {
		slog.Warn("graphx/janusx: ping failed", "address", g.cfg.Address, "error", err)
		recordSpanError(span, err)
		return fmt.Errorf("graphx/janusx: ping: %w", err)
	}
	_, err = result.All()
	if err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("graphx/janusx: ping: %w", err)
	}
	slog.Debug("graphx/janusx: ping ok", "address", g.cfg.Address)
	return nil
}

func recordSpanError(span trace.Span, err error) {
	if err == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// Client returns the underlying *gremlingo.Client for direct use.
func (g *JanusGraph) Client() *gremlingo.Client {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.client
}

// Cfg returns the config used to create this instance.
func (g *JanusGraph) Cfg() graphx.Config { return g.cfg }

// --- singleton ---

var (
	defaultInst = graphx.NewSingleton[JanusGraph]("janusx")
	namedMu     sync.Mutex
	namedInst   = map[string]*JanusGraph{}
)

// Get returns the shared singleton JanusGraph instance, creating it on first call.
// When cfg.Name is non-empty the instance is keyed by name.
func Get(ctx context.Context, cfg graphx.Config) (*JanusGraph, error) {
	if cfg.Name != "" {
		return getNamed(ctx, cfg)
	}
	return defaultInst.Get(func() (*JanusGraph, error) { return New(ctx, cfg) })
}

func getNamed(ctx context.Context, cfg graphx.Config) (*JanusGraph, error) {
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
	namedInst = map[string]*JanusGraph{}
	namedMu.Unlock()
}
