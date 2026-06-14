package agex

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
)

func TestBuild_AGEAddressRequired(t *testing.T) {
	_, err := Build(graphx.Config{})
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestBuild_AGE_SSLModes(t *testing.T) {
	acfg, err := Build(graphx.Config{
		Address:  "localhost:5432",
		Username: "user",
		Password: "pass",
		Database: "testdb",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !contains(acfg.DSN, "sslmode=disable") {
		t.Fatalf("expected sslmode=disable in DSN, got %s", acfg.DSN)
	}

	acfg2, err := Build(graphx.Config{
		Address:  "localhost:5432",
		Username: "user",
		Password: "pass",
		Database: "testdb",
		TLS:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !contains(acfg2.DSN, "sslmode=require") {
		t.Fatalf("expected sslmode=require in DSN, got %s", acfg2.DSN)
	}
}

func TestBuild_AGEDefaultGraph(t *testing.T) {
	acfg, err := Build(graphx.Config{
		Address: "localhost:5432",
	})
	if err != nil {
		t.Fatal(err)
	}
	if acfg.Graph != "graphx" {
		t.Fatalf("expected graph 'graphx', got %s", acfg.Graph)
	}
}

func TestNew_AGENilCtx(t *testing.T) {
	defer Reset()
	_, err := New(nil, graphx.Config{Address: "127.0.0.1:15432"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestGet_AGENamedSingleton(t *testing.T) {
	defer Reset()
	_, err := Get(context.Background(), graphx.Config{Name: "age1", Address: "127.0.0.1:15432"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestReset_AGE(t *testing.T) {
	Reset()
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
