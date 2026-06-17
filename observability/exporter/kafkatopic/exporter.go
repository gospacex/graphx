package kafkatopic

import (
	"context"
	"encoding/json"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type spanExporter struct {
	producer *kafka.Producer
	topic    string
}

func New(producer *kafka.Producer, topic string) trace.SpanExporter {
	return &spanExporter{producer: producer, topic: topic}
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

		headers := map[string]string{}
		otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(headers))

		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &e.topic, Partition: kafka.PartitionAny},
			Value:          payload,
		}
		for k, v := range headers {
			msg.Headers = append(msg.Headers, kafka.Header{Key: k, Value: []byte(v)})
		}
		e.producer.Produce(msg, nil)
	}
	return nil
}

func (e *spanExporter) Shutdown(ctx context.Context) error {
	e.producer.Flush(15 * 1000)
	e.producer.Close()
	return nil
}
