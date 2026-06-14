package janusx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestBuild_JanusAddressRequired(t *testing.T) {
	_, err := Build(graphx.Config{})
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestBuild_JanusDefaultScheme(t *testing.T) {
	jcfg, err := Build(graphx.Config{Address: "localhost:8182"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "ws://127.0.0.1:8182/gremlin"
	if jcfg.URL != expected {
		t.Fatalf("expected %s, got %s", expected, jcfg.URL)
	}
}

func TestBuild_JanusTLSScheme(t *testing.T) {
	jcfg, err := Build(graphx.Config{Address: "127.0.0.1:8182", TLS: true})
	if err != nil {
		t.Fatal(err)
	}
	if jcfg.URL != "wss://127.0.0.1:8182/gremlin" {
		t.Fatalf("expected wss://127.0.0.1:8182/gremlin, got %s", jcfg.URL)
	}
}

func TestNew_JanusNilCtx(t *testing.T) {
	defer Reset()
	_, err := New(nil, graphx.Config{Address: "127.0.0.1:18182"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestGet_JanusNamedSingleton(t *testing.T) {
	defer Reset()
	_, err := Get(context.Background(), graphx.Config{Name: "j1", Address: "127.0.0.1:18182"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestReset_Janus(t *testing.T) {
	Reset()
}
