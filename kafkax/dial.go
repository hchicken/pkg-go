package kafkax

import (
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func (c *KafkaClient) GetDialer() *kafka.Dialer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	if c.opts.sasl != nil {
		dialer.SASLMechanism = c.opts.sasl
	}
	return dialer
}

// dial ...
func (c *KafkaClient) dial() error {
	log.Printf("kafka address: %v ,user: [%v] ", c.opts.address, c.opts.username)

	dialer := c.GetDialer()
	for _, v := range c.opts.address {
		conn, err := dialer.Dial("tcp", v)
		if err != nil {
			log.Printf("Failed to connect to kafka: %v", err)
			return err
		} else {
			log.Println("Connecting to kafka succeeded.")
		}
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
			return err
		}
	}
	return nil
}
