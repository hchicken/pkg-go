package reader

import (
	"crypto/tls"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
)

// ReaderOption defines a function type used for applying options to ReaderOptions.
type ReaderOption func(*ReaderOptions)

// ReaderOptions holds configuration options for the reader.
type ReaderOptions struct {
	reader *kafka.Reader

	sasl    sasl.Mechanism
	address []string

	tls     *tls.Config
	topic   string
	groupId string
}

// newReaderOptions creates a new ReaderOptions instance with the provided options applied.
func newReaderOptions(opts ...ReaderOption) ReaderOptions {
	opt := ReaderOptions{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Topic returns a ReaderOption that sets the topic of ReaderOptions.
func Topic(topic string) ReaderOption {
	return func(o *ReaderOptions) {
		o.topic = topic
	}
}

// GroupId returns a ReaderOption that sets the groupId of ReaderOptions.
func GroupId(id string) ReaderOption {
	return func(o *ReaderOptions) {
		o.groupId = id
	}
}
