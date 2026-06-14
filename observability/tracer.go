package observability

import (
	"context"
	"log/slog"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel"
)

// SetupTracing initializes the global OTel TracerProvider with the exporter and sampler from cfg.
// Returns nil when tracing is disabled or the config is invalid (falls back to noop).
func SetupTracing(ctx context.Context, cfg *Config) error {
	if cfg == nil || !cfg.Enabled {
		slog.Debug("graphx/observability: tracing disabled")
		return nil
	}
	if err := cfg.Validate(); err != nil {
		slog.Warn("graphx/observability: invalid config, using fallback", "error", err)
		return nil
	}

	exporter, err := newSpanExporter(ctx, cfg)
	if err != nil {
		slog.Warn("graphx/observability: exporter init failed, using fallback", "error", err)
		return nil
	}

	sampler := newSampler(cfg)
	res := newResource(cfg.Service)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	slog.Info("graphx/observability: tracing initialized", "exporter", cfg.Exporter, "service", cfg.Service)
	return nil
}

// ShutdownTracerProvider shuts down the global TracerProvider. Safe to call when not initialized.
func ShutdownTracerProvider(ctx context.Context) error {
	tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	if !ok {
		return nil
	}
	if err := tp.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
