package memgraphx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestBuild_MemgraphAddressRequired(t *testing.T) {
	_, err := Build(graphx.Config{})
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestBuild_MemgraphLocalhostRewrite(t *testing.T) {
	nc, err := Build(graphx.Config{Address: "localhost:7687"})
	if err != nil {
		t.Fatal(err)
	}
	if nc.BoltURI != "bolt://127.0.0.1:7687" {
		t.Fatalf("expected bolt://127.0.0.1:7687, got %s", nc.BoltURI)
	}
}

func TestNew_MemgraphNilCtx(t *testing.T) {
	defer Reset()
	_, err := New(nil, graphx.Config{Address: "127.0.0.1:17687"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestGet_MemgraphNamedSingleton(t *testing.T) {
	defer Reset()
	_, err := Get(context.Background(), graphx.Config{Name: "m1", Address: "127.0.0.1:17687"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestReset_Memgraph(t *testing.T) {
	Reset()
}
