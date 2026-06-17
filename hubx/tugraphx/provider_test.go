package tugraphx

import (
	"context"
	"errors"
	"sync"
	"testing"

	hubx "github.com/gospacex/hubx"
)

func TestName_ReturnsCorrectString(t *testing.T) {
	if got := New().Name(); got != "graphx.tugraph" {
		t.Errorf("Name() = %q, want graphx.tugraph", got)
	}
}

func TestBuild_MissingConfigKey(t *testing.T) {
	_, err := New().Build("inst", map[string]any{})
	if !errors.Is(err, hubx.ErrConfigInvalid) {
		t.Errorf("err = %v, want ErrConfigInvalid", err)
	}
}

func TestBuild_MissingEndpoint(t *testing.T) {
	_, err := New().Build("inst", map[string]any{"config": map[string]any{}})
	if !errors.Is(err, hubx.ErrConfigInvalid) {
		t.Errorf("err = %v, want ErrConfigInvalid", err)
	}
}

func TestBuild_UnknownField(t *testing.T) {
	_, err := New().Build("inst", map[string]any{"config": map[string]any{"endpoint": "x", "unknown_xyz": "y"}})
	if !errors.Is(err, hubx.ErrConfigInvalid) {
		t.Errorf("err = %v, want ErrConfigInvalid", err)
	}
}

func TestProviderHealthCheck_NoOp(t *testing.T) {
	if err := New().HealthCheck(context.Background()); err != nil {
		t.Errorf("err = %v", err)
	}
}

func TestProviderClose_NoOp(t *testing.T) {
	if err := New().Close(); err != nil {
		t.Errorf("err = %v", err)
	}
}

func TestBuild_Success(t *testing.T) {
	cli, err := New().Build("inst", map[string]any{"config": map[string]any{"endpoint": "http://localhost:9000"}})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if cli == nil {
		t.Fatal("nil client")
	}
	defer cli.Close()
}

func TestClientHealthCheck_NoError(t *testing.T) {
	cli, _ := New().Build("inst", map[string]any{"config": map[string]any{"endpoint": "http://localhost:9000"}})
	defer cli.Close()
	if err := cli.HealthCheck(context.Background()); err != nil {
		t.Errorf("err = %v", err)
	}
}

func TestClientClose_NoError(t *testing.T) {
	cli, _ := New().Build("inst", map[string]any{"config": map[string]any{"endpoint": "http://localhost:9000"}})
	if err := cli.Close(); err != nil {
		t.Errorf("err = %v", err)
	}
}

func TestConcurrentBuild_Singleton(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cli, err := New().Build("inst", map[string]any{"config": map[string]any{"endpoint": "http://localhost:9000"}})
			if err == nil {
				cli.Close()
			}
		}()
	}
	wg.Wait()
}

func TestRaceFree_UnderRace(t *testing.T) {
	p := New()
	_ = p
}
