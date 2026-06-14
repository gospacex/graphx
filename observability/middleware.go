package observability

import (
	"context"
	"log/slog"
)

// StartQuerySpan starts a child span with graph DB query attributes (db.system, db.statement, etc.).
func StartQuerySpan(ctx context.Context, tr Tracer, backendName, queryLang, query string, params map[string]any) (context.Context, Span) {
	attrs := []any{
		"db.system", backendName,
		"db.statement", query,
	}
	if queryLang != "" {
		attrs = append(attrs, "db.query_language", queryLang)
	}
	if len(params) > 0 {
		attrs = append(attrs, "db.query.parameters", params)
	}
	return tr.Start(ctx, "graphx.query", attrs...)
}

// EndSpan ends a span, logging the error if non-nil.
func EndSpan(s Span, err error) {
	if s == nil {
		return
	}
	if err != nil {
		slog.Debug("graphx/observability: query error", "error", err)
	}
	s.End(err)
}
