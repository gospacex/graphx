// integration_test
//go:build integration

package neo4jx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestIntegration_Neo4jConnect(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, graphx.Config{
		Address:  "localhost:7687",
		Username: "neo4j",
		Password: "password",
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

func TestIntegration_Neo4jSingleton(t *testing.T) {
	defer Reset()
	ctx := context.Background()
	db, err := Get(ctx, graphx.Config{
		Name:     "neo4j-int",
		Address:  "localhost:7687",
		Username: "neo4j",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if db.Driver() == nil {
		t.Fatal("expected non-nil driver")
	}
}
