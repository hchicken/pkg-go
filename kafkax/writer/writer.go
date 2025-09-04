package writer

import (
	"context"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

// Writer represents a Kafka writer.
type Writer struct {
	opts WriterOptions
}

// NewWriter creates and returns a new Writer instance.
func NewWriter(opts ...WriterOption) *Writer {
	options := newOptions(opts...)
	w := &Writer{
		opts: WriterOptions{
			writer: &kafka.Writer{
				Addr:                   kafka.TCP(options.address...),
				Balancer:               options.balancer,
				AllowAutoTopicCreation: options.allowAutoTopicCreation,
			},
		},
	}

	if options.sasl != nil {
		transport := &kafka.Transport{
			SASL: options.sasl,
		}
		if options.tls != nil {
			transport.TLS = options.tls
		}
		w.opts.writer.Transport = transport
	}

	return w
}

// Write sends a message to the specified Kafka topic.
func (w *Writer) Write(ctx context.Context, topic, key, value string) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: []byte(value),
	}

	if err := w.opts.writer.WriteMessages(ctx, msg); err != nil {
		return errors.Wrap(err, "failed to write message")
	}

	return nil
}
