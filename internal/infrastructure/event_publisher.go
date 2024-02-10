package infrastructure

import (
	"context"
	"fmt"

	"github.com/cybre/home-inventory/pkg/domain"
	"github.com/cybre/home-inventory/pkg/kafka"
	"github.com/cybre/home-inventory/pkg/utils"
)

type KafkaProducer interface {
	Produce(context.Context, ...kafka.Record) error
}

type KafkaEventPublisher struct {
	kafkaProducer KafkaProducer
}

func NewKafkaEventPublisher(kafkaProducer KafkaProducer) *KafkaEventPublisher {
	return &KafkaEventPublisher{kafkaProducer: kafkaProducer}
}

func (p *KafkaEventPublisher) PublishEvents(ctx context.Context, events []domain.Event) error {
	records, err := utils.MapWithError(events, func(i uint, event domain.Event) (kafka.Record, error) {
		eventBytes, err := event.Marshal()
		if err != nil {
			return kafka.Record{}, fmt.Errorf("failed to marshal event: %w", err)
		}

		return kafka.Record{
			Key:   event.AggregateID.Marshal(),
			Value: eventBytes,
		}, nil
	})
	if err != nil {
		return fmt.Errorf("failed to map events to kafka records: %w", err)
	}

	if err := p.kafkaProducer.Produce(ctx, records...); err != nil {
		return fmt.Errorf("failed to produce events to kafka: %w", err)
	}

	return nil
}
