package tigergraphx

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/observability"
	"go.opentelemetry.io/otel"
)

func TestTigerGraph_Tracing_Noop(t *testing.T) {
	defer Reset()
	ctx := context.Background()
	cfg := graphx.Config{Address: "127.0.0.1:14240"}
	_, _ = New(ctx, cfg) // MUST NOT panic
}

func TestTigerGraph_Tracing_NewSpanName(t *testing.T) {
	defer Reset()
	tp, exp := observability.NewTestTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	ctx := context.Background()
	_, _ = New(ctx, graphx.Config{Address: "127.0.0.1:14240"})

	for _, s := range exp.GetSpans() {
		if s.Name == "graphx.tigergraphx.new" {
			return
		}
	}
	t.Error("span graphx.tigergraphx.new not found")
}
