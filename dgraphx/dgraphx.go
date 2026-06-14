// Package dgraphx provides a Dgraph backend via HTTP REST API.
package dgraphx

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Dgraph wraps an *http.Client and exposes lifecycle + accessor methods for Dgraph's HTTP API.
type Dgraph struct {
	url    string
	cfg    graphx.Config
	client *http.Client
	tracer trace.Tracer
	mu     sync.RWMutex
}

// New creates a Dgraph HTTP client and returns a ready handle (no pre-connect).
func New(ctx context.Context, cfg graphx.Config) (*Dgraph, error) {
	if ctx == nil {
		slog.Warn("graphx/dgraphx: nil ctx, using context.Background()")
		ctx = context.Background()
	}
	tracer := observability.TracerForBackend("dgraphx", cfg.TracerName)
	ctx, span := tracer.Start(ctx, "graphx.dgraphx.new",
		trace.WithAttributes(
			attribute.String("db.system", "dgraph"),
			attribute.String("server.address", cfg.Address),
			attribute.String("db.operation", "new"),
		))
	defer span.End()

	dcfg, err := Build(cfg)
	if err != nil {
		recordSpanError(span, err)
		return nil, err
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
	}
	if cfg.PoolSize > 0 {
		transport.MaxIdleConns = cfg.PoolSize
		transport.MaxIdleConnsPerHost = cfg.PoolSize
	}
	if cfg.TLS {
		transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return &Dgraph{url: dcfg.BaseURL, cfg: cfg, client: client, tracer: tracer}, nil
}

// Close implements graphx.Closer. Repeated calls are safe.
func (g *Dgraph) Close(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.client == nil {
		return nil
	}
	_, span := g.tracer.Start(ctx, "graphx.dgraphx.close",
		trace.WithAttributes(
			attribute.String("db.system", "dgraph"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "close"),
		))
	defer span.End()
	if transport, ok := g.client.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
	g.client = nil
	return nil
}

// Ping implements graphx.Pinger by sending a POST to /query.
func (g *Dgraph) Ping(ctx context.Context) error {
	g.mu.RLock()
	cl := g.client
	u := g.url
	g.mu.RUnlock()
	if cl == nil {
		return fmt.Errorf("graphx/dgraphx: ping: closed")
	}
	_, span := g.tracer.Start(ctx, "graphx.dgraphx.ping",
		trace.WithAttributes(
			attribute.String("db.system", "dgraph"),
			attribute.String("server.address", g.cfg.Address),
			attribute.String("db.operation", "ping"),
		))
	defer span.End()
	req, err := http.NewRequestWithContext(ctx, "POST", u+"/query", nil)
	if err != nil {
		recordSpanError(span, err)
		return fmt.Errorf("graphx/dgraphx: ping: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cl.Do(req)
	if err != nil {
		slog.Warn("graphx/dgraphx: ping failed", "address", g.cfg.Address, "error", err)
		recordSpanError(span, err)
		return fmt.Errorf("graphx/dgraphx: ping: %w", err)
	}
	resp.Body.Close()
	slog.Debug("graphx/dgraphx: ping ok", "address", g.cfg.Address)
	return nil
}

func recordSpanError(span trace.Span, err error) {
	if err == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// HTTPClient returns the underlying *http.Client for direct use.
func (g *Dgraph) HTTPClient() *http.Client {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.client
}

// Cfg returns the config used to create this instance.
func (g *Dgraph) Cfg() graphx.Config { return g.cfg }

// --- singleton ---

var (
	defaultInst = graphx.NewSingleton[Dgraph]("dgraphx")
	namedMu     sync.Mutex
	namedInst   = map[string]*Dgraph{}
)

// Get returns the shared singleton Dgraph instance, creating it on first call.
// When cfg.Name is non-empty the instance is keyed by name.
func Get(ctx context.Context, cfg graphx.Config) (*Dgraph, error) {
	if cfg.Name != "" {
		return getNamed(ctx, cfg)
	}
	return defaultInst.Get(func() (*Dgraph, error) { return New(ctx, cfg) })
}

func getNamed(ctx context.Context, cfg graphx.Config) (*Dgraph, error) {
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
	namedInst = map[string]*Dgraph{}
	namedMu.Unlock()
}
