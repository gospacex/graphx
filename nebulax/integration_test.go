// integration_test
//go:build integration

package nebulax

import (
	"context"
	"testing"
	"time"

	"github.com/gospacex/graphx"
)

func TestIntegration_NebulaConnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := New(ctx, graphx.Config{
		Address:  "127.0.0.1:9669",
		Username: "root",
		Password: "nebula",
		Database: "graphx",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer db.Close(ctx)

	if db.Pool() == nil {
		t.Fatal("expected non-nil pool")
	}
	if db.Space() != "graphx" {
		t.Fatalf("expected space 'graphx', got %q", db.Space())
	}

	if err := db.Ping(ctx); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestIntegration_NebulaSingleton(t *testing.T) {
	defer Reset()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := Get(ctx, graphx.Config{
		Name:     "nebula-int",
		Address:  "127.0.0.1:9669",
		Username: "root",
		Password: "nebula",
		Database: "graphx",
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if db.Pool() == nil {
		t.Fatal("expected non-nil pool")
	}
	if err := db.Ping(ctx); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestIntegration_NebulaExecute(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := New(ctx, graphx.Config{
		Address:  "127.0.0.1:9669",
		Username: "root",
		Password: "nebula",
		Database: "graphx",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer db.Close(ctx)

	pool := db.Pool()
	result, err := pool.ExecuteAndCheck("SHOW SPACES")
	if err != nil {
		t.Fatalf("SHOW SPACES: %v", err)
	}
	if !result.IsSucceed() {
		t.Fatalf("SHOW SPACES failed: %s", result.GetErrorMsg())
	}
}
