package reader

import (
	"time"

	"github.com/segmentio/kafka-go"
)

// Reader 消费者
type Reader struct {
	opts ReaderOptions
}

// NewReader ...
func NewReader(opts ...ReaderOption) *Reader {
	return newReader(opts...)
}

// newReader ...
func newReader(opts ...ReaderOption) *Reader {

	options := newReaderOptions(opts...)

	// new Reader 对象
	r := new(Reader)
	readerConfig := kafka.ReaderConfig{
		Brokers:        options.address,
		Topic:          options.topic,
		GroupID:        options.groupId,
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // flushes commits to Kafka every second
	}

	// 账号密码设置
	if options.sasl != nil {
		dialer := &kafka.Dialer{
			Timeout:       10 * time.Second,
			DualStack:     true,
			SASLMechanism: options.sasl,
		}
		if options.tls != nil {
			dialer.TLS = options.tls
		}
		readerConfig.Dialer = dialer
	}
	reader := kafka.NewReader(readerConfig)
	r.opts.reader = reader

	return r
}
