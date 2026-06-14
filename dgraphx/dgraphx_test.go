package dgraphx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestBuild_DgraphAddressRequired(t *testing.T) {
	_, err := Build(graphx.Config{})
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestBuild_DgraphPortRewrite(t *testing.T) {
	dcfg, err := Build(graphx.Config{Address: "localhost:9080"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "http://127.0.0.1:8080"
	if dcfg.BaseURL != expected {
		t.Fatalf("expected %s, got %s", expected, dcfg.BaseURL)
	}
}

func TestBuild_DgraphTLSScheme(t *testing.T) {
	dcfg, err := Build(graphx.Config{Address: "127.0.0.1:8080", TLS: true})
	if err != nil {
		t.Fatal(err)
	}
	if dcfg.BaseURL != "https://127.0.0.1:8080" {
		t.Fatalf("expected https://127.0.0.1:8080, got %s", dcfg.BaseURL)
	}
}

func TestNew_DgraphNilCtx(t *testing.T) {
	defer Reset()
	db, err := New(nil, graphx.Config{Address: "127.0.0.1:18080"})
	if err != nil {
		t.Fatal("dgraph client creation should not fail without live connection", err)
	}
	if db.HTTPClient() == nil {
		t.Fatal("expected non-nil HTTP client")
	}
}

func TestGet_DgraphNamedSingleton(t *testing.T) {
	defer Reset()
	db, err := Get(context.Background(), graphx.Config{Name: "d1", Address: "127.0.0.1:18080"})
	if err != nil {
		t.Fatal("dgraph singleton should not fail without live connection", err)
	}
	if db.HTTPClient() == nil {
		t.Fatal("expected non-nil HTTP client")
	}
}

func TestReset_Dgraph(t *testing.T) {
	Reset()
}
