// integration_test
//go:build integration

package agex

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestIntegration_AGEConnect(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, graphx.Config{
		Address:  "localhost:5433",
		Username: "postgres",
		Password: "password",
		Database: "postgres",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer db.Close(ctx)

	if db.DB() == nil {
		t.Fatal("expected non-nil *sql.DB")
	}
	if err := db.Ping(ctx); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}
