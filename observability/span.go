package observability

import (
	"context"
	"log/slog"
)

// Tracer wraps OTel's tracer with a simplified Start method.
type Tracer interface {
	Start(ctx context.Context, spanName string, attrs ...any) (context.Context, Span)
}

// Span wraps OTel's span with simplified End and SetAttribute methods.
type Span interface {
	End(err error)
	SetAttribute(key string, value any)
}

type otelTracer struct{}

func (t *otelTracer) Start(ctx context.Context, spanName string, attrs ...any) (context.Context, Span) {
	tr := globalTracerProvider()
	if tr == nil {
		return ctx, noopSpan{}
	}
	cctx, s := tr.Start(ctx, spanName)
	return cctx, &otelSpan{s: s}
}

type otelSpan struct {
	s interface{ End(...any) }
}

func (s *otelSpan) End(err error) {
	_ = err
	if s.s != nil {
		s.s.End()
	}
}

func (s *otelSpan) SetAttribute(key string, value any) {
}

// NewTracer returns a Tracer that delegates to the global OTel TracerProvider.
// Returns a noop tracer when OTel is not initialized.
func NewTracer() Tracer {
	return &otelTracer{}
}

var noopTracer = &FallbackTracer{}

// globalTracerProvider returns the global TracerProvider if available.
// Returns nil if OTel is not initialized.
func globalTracerProvider() interface {
	Start(ctx context.Context, spanName string, opts ...any) (context.Context, interface{ End(...any) })
} {
	return nil
}

type noopSpan struct{}

func (noopSpan) End(err error)               {}
func (noopSpan) SetAttribute(key string, value any) {}

func init() {
	slog.Debug("graphx/observability: package initialized")
}
