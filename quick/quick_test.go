package quick

import (
	"testing"

	"github.com/gospacex/graphx"
)

func TestN4OS(t *testing.T) {
	_, err := N4OS(graphx.Config{Address: "127.0.0.1:17687"})
	if err == nil {
		t.Fatal("expected connection error for unreachable db")
	}
}

func TestN4PS_InvalidPath(t *testing.T) {
	_, err := N4PS("/nonexistent/graphx.yaml")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestMOS(t *testing.T) {
	_, err := MOS(graphx.Config{Address: "127.0.0.1:17687"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestAOS(t *testing.T) {
	_, err := AOS(graphx.Config{Address: "127.0.0.1:15432"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestGOS(t *testing.T) {
	db, err := GOS(graphx.Config{Address: "127.0.0.1:18080"})
	if err != nil {
		t.Fatal("dgraph quick should not fail without live connection", err)
	}
	if db == nil {
		t.Fatal("expected non-nil db")
	}
}

func TestJOS(t *testing.T) {
	_, err := JOS(graphx.Config{Address: "127.0.0.1:18182"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestVOS(t *testing.T) {
	_, err := VOS(graphx.Config{Address: "127.0.0.1:19669"})
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestTOS(t *testing.T) {
	db, err := TOS(graphx.Config{Address: "127.0.0.1:14240"})
	if err != nil {
		t.Fatal("tigergraph quick should not fail without live connection", err)
	}
	if db == nil {
		t.Fatal("expected non-nil db")
	}
}
