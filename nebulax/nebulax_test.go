package nebulax

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestBuild_NebulaAddressRequired(t *testing.T) {
	_, err := Build(graphx.Config{})
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestBuild_NebulaDefaultPort(t *testing.T) {
	ncfg, err := Build(graphx.Config{Address: "127.0.0.1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ncfg.Hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(ncfg.Hosts))
	}
	if ncfg.Hosts[0] != "127.0.0.1:9669" {
		t.Fatalf("expected 127.0.0.1:9669, got %s", ncfg.Hosts[0])
	}
}

func TestBuild_NebulaMultipleHosts(t *testing.T) {
	ncfg, err := Build(graphx.Config{Address: "127.0.0.1:9669,127.0.0.1:9670"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ncfg.Hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(ncfg.Hosts))
	}
}

func TestBuild_NebulaDefaultSpace(t *testing.T) {
	ncfg, err := Build(graphx.Config{Address: "127.0.0.1:9669"})
	if err != nil {
		t.Fatal(err)
	}
	if ncfg.Space != "default" {
		t.Fatalf("expected space 'default', got %s", ncfg.Space)
	}
}

func TestBuild_NebulaCustomSpace(t *testing.T) {
	ncfg, err := Build(graphx.Config{Address: "127.0.0.1:9669", Database: "mydb"})
	if err != nil {
		t.Fatal(err)
	}
	if ncfg.Space != "mydb" {
		t.Fatalf("expected space 'mydb', got %s", ncfg.Space)
	}
}

func TestNew_NebulaNilCtx(t *testing.T) {
	defer Reset()
	_, err := New(nil, graphx.Config{Address: "127.0.0.1:19669"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestGet_NebulaNamedSingleton(t *testing.T) {
	defer Reset()
	_, err := Get(context.Background(), graphx.Config{Name: "v1", Address: "127.0.0.1:19669"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestReset_Nebula(t *testing.T) {
	Reset()
}

func TestSplitHostPort(t *testing.T) {
	host, port := splitHostPort("127.0.0.1:9669")
	if host != "127.0.0.1" || port != 9669 {
		t.Fatalf("expected 127.0.0.1:9669, got %s:%d", host, port)
	}
}
