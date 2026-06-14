package graphx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_TracerNameFromYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "graphx.yaml")
	yaml := `name: my-app
address: localhost:7687
tracer_name: my-team/graphx
`
	if err := os.WriteFile(path, []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TracerName != "my-team/graphx" {
		t.Errorf("expected TracerName=my-team/graphx, got %q", cfg.TracerName)
	}
}

func TestConfig_TracerNameDefault(t *testing.T) {
	cfg := Config{}
	if cfg.TracerName != "" {
		t.Errorf("expected default TracerName=empty, got %q", cfg.TracerName)
	}
}
