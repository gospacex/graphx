package observability

import (
	"context"
	"log/slog"
)

// FallbackTracer is a slog-based tracer used when OTel is unavailable.
// It logs span start/end events at D-prefixed log levels.
type FallbackTracer struct{}

// NewFallbackTracer returns a FallbackTracer for use when OTel is not initialized.
func NewFallbackTracer() *FallbackTracer {
	return &FallbackTracer{}
}

func (t *FallbackTracer) Start(ctx context.Context, spanName string, attrs ...any) (context.Context, Span) {
	args := make([]any, 0, 2+len(attrs))
	args = append(args, "span", spanName)
	args = append(args, attrs...)
	slog.Debug("graphx/observability: start", args...)
	return ctx, &fallbackSpan{name: spanName}
}

type fallbackSpan struct {
	name string
	err  error
}

func (s *fallbackSpan) End(err error) {
	s.err = err
	if err != nil {
		slog.Warn("graphx/observability: end", "span", s.name, "error", err)
	} else {
		slog.Debug("graphx/observability: end", "span", s.name)
	}
}

func (s *fallbackSpan) SetAttribute(key string, value any) {
	slog.Debug("graphx/observability: attr", "span", s.name, key, value)
}
