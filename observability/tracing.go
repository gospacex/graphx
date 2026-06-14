package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// InjectTrace serialises the trace context from ctx into a string map for propagation.
func InjectTrace(ctx context.Context) (map[string]string, bool) {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	if len(carrier) == 0 {
		return nil, false
	}
	return map[string]string(carrier), true
}

// ExtractTrace deserialises a string map into a context with trace context.
func ExtractTrace(ctx context.Context, carrier map[string]string) context.Context {
	if len(carrier) == 0 {
		return ctx
	}
	mc := propagation.MapCarrier(carrier)
	return otel.GetTextMapPropagator().Extract(ctx, mc)
}
