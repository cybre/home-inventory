package infrastructure

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/cybre/home-inventory/pkg/eventsourcing"
	"github.com/cybre/home-inventory/pkg/kafka"
	"github.com/cybre/home-inventory/pkg/logging"
	"github.com/cybre/home-inventory/pkg/utils"
)

type EventHandler interface {
	HandleEvent(ctx context.Context, event eventsourcing.EventData) error
	Events() []eventsourcing.EventType
	Name() string
}

type KafkaEventMessaging struct {
	logger        *slog.Logger
	producer      *kafka.Producer
	brokers       []string
	topic         string
	eventHandlers []EventHandler
	consumers     []*kafka.Consumer
}

func NewKafkaEventMessaging(brokers []string, topic string, logger *slog.Logger) (*KafkaEventMessaging, error) {
	producer, err := kafka.NewProducer(brokers, topic, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &KafkaEventMessaging{
		logger:        logger,
		producer:      producer,
		brokers:       brokers,
		topic:         topic,
		eventHandlers: []EventHandler{},
		consumers:     []*kafka.Consumer{},
	}, nil
}

func (p *KafkaEventMessaging) PublishEvents(ctx context.Context, events []eventsourcing.Event) error {
	records, err := utils.MapWithError(events, func(i uint, event eventsourcing.Event) (kafka.Record, error) {
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

	if err := p.producer.Produce(ctx, records...); err != nil {
		return fmt.Errorf("failed to produce events to kafka: %w", err)
	}

	return nil
}

func (c *KafkaEventMessaging) ConsumeEvents(ctx context.Context, handler EventHandler) error {
	kafkaConsumer, err := kafka.NewConsumer(
		c.brokers,
		c.topic,
		handler.Name(),
		c.logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	c.consumers = append(c.consumers, kafkaConsumer)

	go kafkaConsumer.Consume(ctx, func(record kafka.Record) {
		event, err := eventsourcing.UnmarshalEvent(record.Value)
		if err != nil {
			c.logger.Error("failed to unmarshal event", slog.Any("error", err))
			return
		}

		if !slices.Contains(handler.Events(), event.EventType) {
			return
		}

		handlerContext := logging.WithLogger(
			ctx,
			logging.FromContext(ctx).With(
				slog.String("event_handler", handler.Name()),
				slog.Any("event_type", event.EventType),
			),
		)
		if err := handler.HandleEvent(handlerContext, event.Data); err != nil {
			c.logger.Error("failed to handle event", slog.Any("error", err))
		}
	})

	return nil
}

func (c *KafkaEventMessaging) Close() {
	for _, consumer := range c.consumers {
		consumer.Close()
	}
	c.producer.Close()
}
