//go:build integration

package agex

import (
	"context"
	"os"
	"testing"
)

func TestIntegration_BuildAndPing(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") != "" {
		t.Skip()
	}
	ep := os.Getenv("AGE_ENDPOINT")
	if ep == "" {
		t.Skip()
	}
	cli, err := New().Build("it", map[string]any{"config": map[string]any{"endpoint": ep}})
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	if err := cli.HealthCheck(context.Background()); err != nil {
		t.Fatal(err)
	}
}