// integration_test
//go:build integration

package janusx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestIntegration_JanusConnect(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, graphx.Config{
		Address: "127.0.0.1:8182",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer db.Close(ctx)

	if db.Client() == nil {
		t.Fatal("expected non-nil Gremlin client")
	}
}
