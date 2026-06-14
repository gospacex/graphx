package agex

import (
	"context"
	"testing"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/observability"
	"go.opentelemetry.io/otel"
)

func TestAgex_Tracing_Noop(t *testing.T) {
	defer Reset()
	ctx := context.Background()
	cfg := graphx.Config{Address: "127.0.0.1:5432", Username: "u", Password: "p"}
	_, _ = New(ctx, cfg) // MUST NOT panic
}

func TestAgex_Tracing_NewSpanName(t *testing.T) {
	defer Reset()
	tp, exp := observability.NewTestTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	ctx := context.Background()
	_, _ = New(ctx, graphx.Config{Address: "127.0.0.1:5432", Username: "u", Password: "p"})

	for _, s := range exp.GetSpans() {
		if s.Name == "graphx.agex.new" {
			return
		}
	}
	t.Error("span graphx.agex.new not found")
}
