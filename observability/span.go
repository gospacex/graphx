package observability

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	tp := otel.GetTracerProvider()
	tp, ok := tp.(trace.TracerProvider)
	if !ok {
		return ctx, noopSpan{}
	}
	tracer := tp.Tracer("graphx")
	cctx, s := tracer.Start(ctx, spanName)
	os := &otelSpan{s: s}
	for i := 0; i < len(attrs)-1; i += 2 {
		if key, ok := attrs[i].(string); ok {
			os.SetAttribute(key, attrs[i+1])
		}
	}
	return cctx, os
}

type otelSpan struct {
	s trace.Span
}

func (s *otelSpan) End(err error) {
	if s.s == nil {
		return
	}
	if err != nil {
		s.s.RecordError(err)
		s.s.SetStatus(codes.Error, err.Error())
	}
	s.s.End()
}

func (s *otelSpan) SetAttribute(key string, value any) {
	if s.s == nil {
		return
	}
	s.s.SetAttributes(attribute.String(key, fmt.Sprintf("%v", value)))
}

// NewTracer returns a Tracer that delegates to the global OTel TracerProvider.
// Returns a noop tracer when OTel is not initialized.
func NewTracer() Tracer {
	return &otelTracer{}
}

var noopTracer = &FallbackTracer{}

type noopSpan struct{}

func (noopSpan) End(err error)               {}
func (noopSpan) SetAttribute(key string, value any) {}

func init() {
	slog.Debug("graphx/observability: package initialized")
}
