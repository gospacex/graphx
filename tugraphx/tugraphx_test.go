package tugraphx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestBuild_TuGraphAddressRequired(t *testing.T) {
	_, err := Build(graphx.Config{})
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestBuild_TuGraphDefaultScheme(t *testing.T) {
	tcfg, err := Build(graphx.Config{Address: "localhost:7688"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "bolt://127.0.0.1:7688"
	if tcfg.BoltURI != expected {
		t.Fatalf("expected %s, got %s", expected, tcfg.BoltURI)
	}
}

func TestBuild_TuGraphTLSScheme(t *testing.T) {
	tcfg, err := Build(graphx.Config{Address: "127.0.0.1:7688", TLS: true})
	if err != nil {
		t.Fatal(err)
	}
	if tcfg.BoltURI != "bolts://127.0.0.1:7688" {
		t.Fatalf("expected bolts://127.0.0.1:7688, got %s", tcfg.BoltURI)
	}
}

func TestNew_TuGraphNilCtx(t *testing.T) {
	defer Reset()
	_, err := New(nil, graphx.Config{Address: "127.0.0.1:17688"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestGet_TuGraphNamedSingleton(t *testing.T) {
	defer Reset()
	_, err := Get(context.Background(), graphx.Config{Name: "t1", Address: "127.0.0.1:17688"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestReset_TuGraph(t *testing.T) {
	Reset()
}
