package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client *kgo.Client
}

func NewConsumer(brokers []string, topic, consumerGroup string) (*Consumer, error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(consumerGroup),
		kgo.ConsumeTopics(topic),
	)

	if err != nil {
		return nil, fmt.Errorf("error initializing Kafka consumer: %w", err)
	}

	return &Consumer{client: cl}, nil
}

func (c *Consumer) Consume(ctx context.Context, callback func(Record)) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fetches := c.client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return fmt.Errorf("Kafka client closed")
		}

		if fetches.Err() != nil {
			return fmt.Errorf("error consuming message from Kafka: %w", fetches.Err())
		}

		fetches.EachRecord(func(record *kgo.Record) {
			callback(Record{
				Key:   record.Key,
				Value: record.Value,
			})
		})

		time.Sleep(1 * time.Second)
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}
