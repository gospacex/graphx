// integration_test
//go:build integration

package memgraphx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestIntegration_MemgraphConnect(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, graphx.Config{
		Address: "127.0.0.1:7689",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer db.Close(ctx)

	if db.Driver() == nil {
		t.Fatal("expected non-nil driver")
	}
	if err := db.Ping(ctx); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}
