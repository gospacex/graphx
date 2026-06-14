package observability

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// TracerForBackend returns an OTel tracer for the given backend subpackage,
// honoring an optional cfgTracerName override. When cfgTracerName is empty
// the convention "graphx/<backend>" is used (e.g., "graphx/neo4jx").
//
// The returned tracer is the global OTel TracerProvider's named tracer.
// If OTel is not initialized (no SetupTracing call), the returned tracer
// is the OTel noop tracer — Start/End are zero-cost.
func TracerForBackend(backendName, cfgTracerName string) trace.Tracer {
	name := cfgTracerName
	if name == "" {
		name = "graphx/" + backendName
	}
	return otel.Tracer(name)
}
