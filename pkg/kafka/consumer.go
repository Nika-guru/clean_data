package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// StartOffset enum
const (
	ComsumeAllData     = kafka.FirstOffset //Consume from beginning
	ConsumeNewDataOnly = kafka.LastOffset  //Consume new data only after Comsumer start listening
)

// Local Cache Configuration Struct
type ConsumerConfig struct {
	Brokers []string
	Topic   string

	/*
		GroupID holds the optional consumer group id.
		If GroupID is specified, then Partition should NOT be specified.
	*/
	GroupId string

	/*
		MinBytes - If the consumer polls the cluster to check if there is any new data on the topic for the my-group consumer ID, the cluster will only respond if there are at least 5 new bytes of information to send.
	*/
	MinBytes int // default = 1 Bytes

	// The kafka-go library requires you to set the MaxBytes in case the MinBytes are set
	MaxBytes int // default = 1e6 Bytes

	// Wait for at most MaxWait seconds before receiving new data
	MaxWait time.Duration // default = 10 seconds

	/*
		StartOffset this only applies for new consumer groups. If you’ve already consumed data with the same consumer GroupID setting before, you will continue from wherever you left off.
	*/
	StartOffset int64 // Default ComsumeAllData
}

type Consumer struct {
	Context context.Context
	Reader  *kafka.Reader
}

// Consumer Create Method
func CreateConsumer(config ConsumerConfig) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Brokers,
		Topic:       config.Topic,
		GroupID:     config.GroupId,
		MinBytes:    config.MinBytes,
		MaxBytes:    config.MaxBytes,
		MaxWait:     config.MaxWait,
		StartOffset: config.StartOffset,
	})
	return &Consumer{
		Context: context.Background(),
		Reader:  reader,
	}
}

func (consumer *Consumer) Consume() ([]byte, error) {
	data, err := consumer.Reader.ReadMessage(consumer.Context)
	return data.Value, err
}
