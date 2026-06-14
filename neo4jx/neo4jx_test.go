package neo4jx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestBuild_AddressRequired(t *testing.T) {
	_, err := Build(graphx.Config{})
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestBuild_DefaultScheme(t *testing.T) {
	nc, err := Build(graphx.Config{Address: "localhost:7687"})
	if err != nil {
		t.Fatal(err)
	}
	if nc.BoltURI != "bolt://localhost:7687" {
		t.Fatalf("expected bolt://localhost:7687, got %s", nc.BoltURI)
	}
}

func TestBuild_CustomScheme(t *testing.T) {
	nc, err := Build(graphx.Config{
		Address: "localhost:7687",
		Extra:   map[string]any{"scheme": "neo4j"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if nc.BoltURI != "neo4j://localhost:7687" {
		t.Fatalf("expected neo4j://localhost:7687, got %s", nc.BoltURI)
	}
}

func TestNew_NilCtx(t *testing.T) {
	defer Reset()
	_, err := New(nil, graphx.Config{Address: "127.0.0.1:17687"})
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
}

func TestGet_DefaultSingleton(t *testing.T) {
	defer Reset()
	// With no real server, Get should try to connect and fail
	_, err := Get(context.Background(), graphx.Config{Address: "127.0.0.1:17687"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestGet_NamedSingleton(t *testing.T) {
	defer Reset()
	_, err := Get(context.Background(), graphx.Config{Name: "test", Address: "127.0.0.1:17687"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestReset(t *testing.T) {
	if err := recover(); err != nil {
		t.Fatal("Reset should not panic")
	}
	Reset()
}
