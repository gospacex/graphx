package observability

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// NewTestTracerProvider returns a TracerProvider paired with an in-memory
// exporter, suitable for unit tests that need to assert on emitted spans
// without depending on a live OTel collector.
//
// Usage:
//
//	tp, exp := observability.NewTestTracerProvider()
//	otel.SetTracerProvider(tp)
//	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })
//	// ... exercise code that calls otel.Tracer(...).Start(...)
//	spans := exp.GetSpans()
//
// The returned TracerProvider and exporter are safe for concurrent use.
func NewTestTracerProvider() (*sdktrace.TracerProvider, *tracetest.InMemoryExporter) {
	exp := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	return tp, exp
}
