package neo4jx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// TestNeo4j_Tracing_Noop asserts backend operations don't panic when
// observability is NOT initialized.
func TestNeo4j_Tracing_Noop(t *testing.T) {
	defer Reset()
	ctx := context.Background()
	cfg := graphx.Config{Address: "127.0.0.1:7687", Username: "neo4j", Password: "p"}
	_, _ = New(ctx, cfg)
}

// TestNeo4j_Tracing_NewSpanName verifies the span name emitted by New.
func TestNeo4j_Tracing_NewSpanName(t *testing.T) {
	defer Reset()

	tp, exp := observability.NewTestTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	ctx := context.Background()
	cfg := graphx.Config{Address: "127.0.0.1:7687", Username: "neo4j", Password: "p"}
	_, _ = New(ctx, cfg)

	var found bool
	for _, s := range exp.GetSpans() {
		if s.Name == "graphx.neo4jx.new" {
			found = true
			break
		}
	}
	if !found {
		var names []string
		for _, s := range exp.GetSpans() {
			names = append(names, s.Name)
		}
		t.Errorf("expected span 'graphx.neo4jx.new', got: %v", names)
	}
}

// TestNeo4j_Tracing_NewSpanAttributes verifies the span attributes.
func TestNeo4j_Tracing_NewSpanAttributes(t *testing.T) {
	defer Reset()

	tp, exp := observability.NewTestTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	ctx := context.Background()
	cfg := graphx.Config{Address: "127.0.0.1:7687", Username: "neo4j", Password: "p"}
	_, _ = New(ctx, cfg)

	for _, s := range exp.GetSpans() {
		if s.Name != "graphx.neo4jx.new" {
			continue
		}
		want := map[attribute.Key]string{
			"db.system":      "neo4j",
			"server.address": "127.0.0.1:7687",
			"db.operation":   "new",
		}
		got := map[attribute.Key]attribute.Value{}
		for _, a := range s.Attributes {
			got[a.Key] = a.Value
		}
		for k, v := range want {
			if gv, ok := got[k]; !ok {
				t.Errorf("missing attribute %s", k)
			} else if gv.AsString() != v {
				t.Errorf("attribute %s: want %s, got %s", k, v, gv.AsString())
			}
		}
		return
	}
	t.Error("span graphx.neo4jx.new not found")
}
