// integration_test
//go:build integration

package tugraphx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestIntegration_TuGraphConnect(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, graphx.Config{
		Address:  "localhost:7688",
		Username: "admin",
		Password: "73@TuGraph",
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

func TestIntegration_TuGraphSingleton(t *testing.T) {
	defer Reset()
	ctx := context.Background()
	db, err := Get(ctx, graphx.Config{
		Name:     "tugraph-int",
		Address:  "localhost:7688",
		Username: "admin",
		Password: "73@TuGraph",
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if db.Driver() == nil {
		t.Fatal("expected non-nil driver")
	}
}
