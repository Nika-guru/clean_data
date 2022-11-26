package kafka

import (
	"context"
	"encoding/json"
	"errors"
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

type Consumer[T EventData] struct {
	context    context.Context
	reader     *kafka.Reader
	eventTypes []string
}

// Consumer Create Method
func CreateConsumer(config ConsumerConfig, eventTypes []string) *Consumer[EventData] {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Brokers,
		Topic:       config.Topic,
		GroupID:     config.GroupId,
		MinBytes:    config.MinBytes,
		MaxBytes:    config.MaxBytes,
		MaxWait:     config.MaxWait,
		StartOffset: config.StartOffset,
	})
	return &Consumer[EventData]{
		context:    context.Background(),
		reader:     reader,
		eventTypes: eventTypes,
	}
}

func (consumer *Consumer[T]) Consume() (*Schema[T], error) {
	data, err := consumer.reader.ReadMessage(consumer.context)
	if err != nil {
		return nil, err
	}

	var schema = &Schema[T]{}
	err = json.Unmarshal(data.Value, &schema)
	if err != nil {
		return nil, err
	}

	if !schema.ContainsEventTypes(consumer.eventTypes) {
		return nil, errors.New("imvalid event type")
	}
	return schema, err
}
