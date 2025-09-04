package writer

import (
	"crypto/tls"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
)

// WriterOption defines a function to configure WriterOptions
type WriterOption func(*WriterOptions)

// WriterOptions holds configuration for the Kafka writer
type WriterOptions struct {
	writer                 *kafka.Writer
	sasl                   sasl.Mechanism
	address                []string
	tls                    *tls.Config
	balancer               kafka.Balancer
	allowAutoTopicCreation bool
}

// newOptions creates a new WriterOptions with default values and applies given options
func newOptions(opts ...WriterOption) WriterOptions {
	opt := WriterOptions{
		balancer:               &kafka.LeastBytes{},
		allowAutoTopicCreation: true,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Address sets the Kafka broker addresses
func Address(address ...string) WriterOption {
	return func(o *WriterOptions) {
		o.address = address
	}
}

// SASL sets the SASL mechanism for authentication
func SASL(mechanism sasl.Mechanism) WriterOption {
	return func(o *WriterOptions) {
		o.sasl = mechanism
	}
}

// Balancer sets the balancer for distributing messages
func Balancer(balancer kafka.Balancer) WriterOption {
	return func(o *WriterOptions) {
		o.balancer = balancer
	}
}

// AllowAutoTopicCreation sets whether to allow auto topic creation
func AllowAutoTopicCreation(allow bool) WriterOption {
	return func(o *WriterOptions) {
		o.allowAutoTopicCreation = allow
	}
}

// TLS sets the TLS configuration
func TLS(config *tls.Config) WriterOption {
	return func(o *WriterOptions) {
		o.tls = config
	}
}
