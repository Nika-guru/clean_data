package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// RequiredAcks enum
const (
	AcksRequireNone = 0  // fire-and-forget, do not wait for acknowledgements from the
	AcksRequireOne  = 1  // only wait for the leader to acknowledge - good for non-transactional data
	AcksRequireAll  = -1 // wait for all brokers to acknowledge the writes (In-Sync Replicas)
)

// Local Cache Configuration Struct
type ProducerConfig struct {
	Brokers []string
	Topic   string

	//BatchSize - The total number of messages that should be buffered before writing to the brokers
	BatchSize int //Default value is 100

	/*
		BatchTimeout - The maximum time before which messages are written to the brokers. That means that even if the message batch is not full, they will still be written onto the Kafka cluster once this time period has elapsed.
	*/
	BatchTimeout time.Duration //Default value is 1 seconds

	/*
		RequiredAcks have 3 options
		1. All brokers acknowledge that they have received the message (-1)
		2. Only the leading broker acknowledges that it has received the messages (1). The remaining brokers can still eventually receive the message, but we won’t wait for them to do so.
		3. No one acknowledges receiving the message (0). This is basically a fire-and-forget mode, where we don’t care if our message is received or not. This should only be used for data that you are ok with losing a bit of, but require high throughput for.

	*/
	RequiredAcks int //Default value is AcksRequireOne(1)
}

type Producer struct {
	Context context.Context
	Writer  *kafka.Writer
}

// Producer Create Method
func CreateProducer(config ProducerConfig) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      config.Brokers,
		Topic:        config.Topic,
		BatchSize:    config.BatchSize,
		BatchTimeout: config.BatchTimeout,
		RequiredAcks: config.RequiredAcks,
	})
	return &Producer{
		Context: context.Background(),
		Writer:  writer,
	}
}

func (producer *Producer) Produce(data []byte) error {
	return producer.Writer.WriteMessages(producer.Context, kafka.Message{
		Value: data,
	})
}
