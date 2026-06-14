package observability

import (
	"context"
	"sync"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// TestNewTestTracerProvider_CapturesSpans verifies that spans created via
// the global TracerProvider are captured by the in-memory exporter.
func TestNewTestTracerProvider_CapturesSpans(t *testing.T) {
	tp, exp := NewTestTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	_, span := tp.Tracer("test").Start(context.Background(), "test.span")
	span.SetAttributes(attribute.String("key", "value"))
	span.End()

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if spans[0].Name != "test.span" {
		t.Errorf("expected name=test.span, got %s", spans[0].Name)
	}
}

// TestNewTestTracerProvider_ConcurrentSafe verifies the in-memory exporter
// is safe under concurrent span emission (run under -race).
func TestNewTestTracerProvider_ConcurrentSafe(t *testing.T) {
	tp, exp := NewTestTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })

	const goroutines = 16
	const spansPerG = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			tr := tp.Tracer("concurrent")
			for j := 0; j < spansPerG; j++ {
				_, span := tr.Start(context.Background(), "op")
				span.End()
			}
		}()
	}
	wg.Wait()

	all := exp.GetSpans()
	if len(all) != goroutines*spansPerG {
		t.Errorf("expected %d spans, got %d", goroutines*spansPerG, len(all))
	}
}
