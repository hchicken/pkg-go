package kafkax

import (
	"kafkax/writer"
)

// KafkaClient ...
type KafkaClient struct {
	opts Options
}

func NewKafkaClient(opts ...Option) (*KafkaClient, error) {
	cli := new(KafkaClient)

	options := newOptions(opts...)
	cli.opts = options

	// 连通性判断
	if err := cli.dial(); err != nil {
		return nil, err
	}

	return cli, nil
}

// NewReader ...
//func (*KafkaClient) NewReader(...writer.WriterOptions) reader.Reader {
//	writer.NewWriter()
//}

// NewWriter ...
func (client *KafkaClient) NewWriter(opts ...writer.WriterOption) *writer.Writer {
	opts = append(opts,
		writer.SASL(client.opts.sasl),
		writer.Address(client.opts.address...),
	)
	w := writer.NewWriter(opts...)
	return w
}
