// integration_test
//go:build integration

package dgraphx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestIntegration_DgraphConnect(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, graphx.Config{
		Address: "localhost:9080",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer db.Close(ctx)

	if db.HTTPClient() == nil {
		t.Fatal("expected non-nil HTTP client")
	}
}
