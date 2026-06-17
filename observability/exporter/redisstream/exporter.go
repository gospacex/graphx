package redisstream

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/sdk/trace"
)

type spanExporter struct {
	client *redis.Client
	stream string
}

func New(client *redis.Client, stream string) trace.SpanExporter {
	return &spanExporter{client: client, stream: stream}
}

type spanData struct {
	TraceID    string            `json:"trace_id"`
	SpanID     string            `json:"span_id"`
	Name       string            `json:"name"`
	StartTime  time.Time         `json:"start_time"`
	EndTime    time.Time         `json:"end_time"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

func (e *spanExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	for _, s := range spans {
		attr := map[string]string{}
		for _, a := range s.Attributes() {
			attr[string(a.Key)] = a.Value.AsString()
		}
		data := spanData{
			TraceID:   s.SpanContext().TraceID().String(),
			SpanID:    s.SpanContext().SpanID().String(),
			Name:      s.Name(),
			StartTime: s.StartTime(),
			EndTime:   s.EndTime(),
		}
		if len(attr) > 0 {
			data.Attributes = attr
		}
		payload, _ := json.Marshal(data)

		args := &redis.XAddArgs{
			Stream: e.stream,
			Values: map[string]any{"span": string(payload)},
			MaxLen: 1000,
			ID:     "*",
		}
		if err := e.client.XAdd(ctx, args).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (e *spanExporter) Shutdown(ctx context.Context) error {
	return e.client.Close()
}
