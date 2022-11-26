package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"review-service/pkg/log"
	"time"

	"github.com/segmentio/kafka-go"
)

// RequiredAcks enum
const (
	// fire-and-forget, do not wait for acknowledgements from the cluster
	AcksRequireNone = kafka.RequireNone
	// only wait for the leader to acknowledge - good for non-transactional data
	AcksRequireOne = kafka.RequireOne
	// wait for all brokers to acknowledge the writes (In-Sync Replicas)
	AcksRequireAll = kafka.RequireAll
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
	RequiredAcks kafka.RequiredAcks //Default value is AcksRequireOne(1)
}

const (
	_defaultKafkaTimeout = 250 * time.Millisecond
)

type Producer[T EventData] struct {
	writer     *kafka.Writer
	eventTypes []string
}

// Producer Create Method
func CreateProducer(config ProducerConfig, eventTypes []string) *Producer[EventData] {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers...),
		Topic:                  config.Topic,
		BatchSize:              config.BatchSize,
		BatchTimeout:           config.BatchTimeout,
		RequiredAcks:           config.RequiredAcks,
		AllowAutoTopicCreation: true,
	}
	return &Producer[EventData]{
		writer:     writer,
		eventTypes: eventTypes,
	}
}

func (producer *Producer[T]) Produce(schema *Schema[T]) error {
	if !schema.ContainsEventTypes(producer.eventTypes) {
		return errors.New("invalid event type")
	}

	jsonData, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	const retries = 3
	for i := 0; i < retries; i++ {
		// attempt to create topic prior to publishing the message
		err = producer.writer.WriteMessages(context.Background(), kafka.Message{
			Value: jsonData,
		})

		//LeaderNotAvailable | Message send timeout or errors
		if err != nil {
			log.Println(log.LogLevelError, "kafka-write: ", err)
			fmt.Println("Kafka wait ...")
			time.Sleep(_defaultKafkaTimeout)
			continue
		} else {
			break
		}
	}
	return nil
}
