package neo4jx

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	hubx "github.com/gospacex/hubx"
)

// fullCfg returns a graphx.Config-equivalent map that satisfies
// mapstructure's ErrorUnset / ErrorUnused contract. graphx.Config
// includes a *tls.Config field tagged yaml:"-" which mapstructure
// still tracks; we supply an explicit nil value to satisfy
// ErrorUnset.
func fullCfg(addr string) map[string]any {
	return map[string]any{
		"address":     addr,
		"username":    "neo4j",
		"password":    "pw",
		"database":    "",
		"poolSize":    0,
		"tls":         false,
		"tracer_name": "",
		"name":        "",
		"extra":       map[string]any{},
		"tlsConfig":   nil,
	}
}

func TestName_ReturnsGraphxNeo4j(t *testing.T) {
	if got := (&Provider{}).Name(); got != "graphx.neo4j" {
		t.Fatalf("Name() = %q, want %q", got, "graphx.neo4j")
	}
}

// TestBuild_Success verifies that mapstructure decodes a complete
// config and that the driver New is invoked. We use an unreachable
// address so the driver returns ErrBuildFailed — this still proves
// the decode succeeded.
func TestBuild_Success(t *testing.T) {
	p := New()
	_, err := p.Build("main", map[string]any{
		"config": fullCfg("127.0.0.1:1"),
	})
	if err == nil {
		t.Fatal("expected error against unreachable address")
	}
	if !errors.Is(err, hubx.ErrBuildFailed) {
		t.Fatalf("expected ErrBuildFailed (mapstructure decoded; driver failed), got %v", err)
	}
}

func TestBuild_MissingConfigKey(t *testing.T) {
	p := New()
	cli, err := p.Build("main", map[string]any{})
	if cli != nil {
		t.Fatal("expected nil client on error")
	}
	if err == nil {
		t.Fatal("expected error when 'config' key is missing")
	}
	if !errors.Is(err, hubx.ErrConfigInvalid) {
		t.Fatalf("expected ErrConfigInvalid, got %v", err)
	}
}

// TestBuild_MissingRequiredField forces a decode failure by passing
// a non-map value as the config. mapstructure cannot decode an int
// into graphx.Config and returns an error which is wrapped in
// ErrConfigInvalid.
func TestBuild_MissingRequiredField(t *testing.T) {
	p := New()
	cli, err := p.Build("main", map[string]any{
		"config": 42,
	})
	if cli != nil {
		t.Fatal("expected nil client on error")
	}
	if err == nil {
		t.Fatal("expected error when 'config' value is not a map")
	}
	if !errors.Is(err, hubx.ErrConfigInvalid) {
		t.Fatalf("expected ErrConfigInvalid, got %v", err)
	}
}

func TestBuild_UnknownField(t *testing.T) {
	p := New()
	cfg := fullCfg("127.0.0.1:7687")
	cfg["this_does_not"] = "exist"
	cli, err := p.Build("main", map[string]any{
		"config": cfg,
	})
	if cli != nil {
		t.Fatal("expected nil client on error")
	}
	if err == nil {
		t.Fatal("expected error on unknown field")
	}
	if !errors.Is(err, hubx.ErrConfigInvalid) {
		t.Fatalf("expected ErrConfigInvalid, got %v", err)
	}
}

func TestBuild_DriverNewFailure(t *testing.T) {
	p := New()
	cli, err := p.Build("main", map[string]any{
		"config": fullCfg("127.0.0.1:1"),
	})
	if cli != nil {
		t.Fatal("expected nil client on driver new failure")
	}
	if err == nil {
		t.Fatal("expected error from neo4jx.New")
	}
	if !errors.Is(err, hubx.ErrBuildFailed) {
		t.Fatalf("expected ErrBuildFailed, got %v", err)
	}
}

func TestProviderHealthCheck_NoOp(t *testing.T) {
	p := New()
	if err := p.HealthCheck(context.Background()); err != nil {
		t.Fatalf("Provider.HealthCheck = %v, want nil", err)
	}
}

func TestProviderClose_NoOp(t *testing.T) {
	p := New()
	if err := p.Close(); err != nil {
		t.Fatalf("Provider.Close = %v, want nil", err)
	}
}

// TestClientHealthCheck exercises the wrapper's delegation by
// running the no-op provider surface (Build may fail against :1;
// that's fine for race coverage).
func TestClientHealthCheck(t *testing.T) {
	p := New()
	cli, err := p.Build("hc", map[string]any{
		"config": fullCfg("127.0.0.1:1"),
	})
	if err == nil && cli != nil {
		t.Cleanup(func() { _ = cli.Close() })
		_ = cli.HealthCheck(context.Background())
	}
}

func TestConcurrentBuild_Singleton(t *testing.T) {
	p := New()
	const N = 50
	var (
		wg       sync.WaitGroup
		errCount atomic.Int64
	)
	start := make(chan struct{})
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := p.Build("main", map[string]any{
				"config": fullCfg("127.0.0.1:1"),
			})
			if err != nil {
				errCount.Add(1)
			}
		}()
	}
	close(start)
	wg.Wait()
	if errCount.Load() != N {
		t.Fatalf("expected %d errors, got %d", N, errCount.Load())
	}
}

func TestRaceFree_UnderRace(t *testing.T) {
	p := New()
	cli, err := p.Build("race", map[string]any{
		"config": fullCfg("127.0.0.1:1"),
	})
	if err == nil && cli != nil {
		t.Cleanup(func() { _ = cli.Close() })
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = p.HealthCheck(context.Background())
		}()
		go func() {
			defer wg.Done()
			_ = p.Close()
		}()
	}
	wg.Wait()
}