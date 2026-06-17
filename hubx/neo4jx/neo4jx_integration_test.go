//go:build integration

package neo4jx

import (
	"context"
	"os"
	"testing"
)

func TestIntegration_BuildAndPing(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") != "" {
		t.Skip("SKIP_INTEGRATION set")
	}
	host := os.Getenv("NEO4J_HOST")
	if host == "" {
		host = "127.0.0.1:7687"
	}
	user := os.Getenv("NEO4J_USER")
	if user == "" {
		user = "neo4j"
	}
	pass := os.Getenv("NEO4J_PASS")
	if pass == "" {
		pass = "password"
	}

	p := New()
	cli, err := p.Build("main", map[string]any{
		"config": map[string]any{
			"address":  host,
			"username": user,
			"password": pass,
		},
	})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	t.Cleanup(func() { _ = cli.Close() })

	if err := cli.HealthCheck(context.Background()); err != nil {
		t.Fatalf("HealthCheck: %v", err)
	}
	if err := p.HealthCheck(context.Background()); err != nil {
		t.Fatalf("Provider.HealthCheck: %v", err)
	}
}