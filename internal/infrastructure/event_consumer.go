package infrastructure

import (
	"context"
	"fmt"
	"slices"

	"github.com/cybre/home-inventory/internal/infrastructure/kafka"
	"github.com/cybre/home-inventory/pkg/domain"
)

type EventHandler interface {
	HandleEvent(ctx context.Context, event domain.EventData) error
	Events() []domain.EventType
	Name() string
}

type KafkaEventConsumer struct {
	eventHandlers []EventHandler
	consumers     []*kafka.Consumer
	brokers       []string
	topic         string
}

func NewKafkaEventConsumer(brokers []string, topic string) *KafkaEventConsumer {
	return &KafkaEventConsumer{
		eventHandlers: []EventHandler{},
		consumers:     []*kafka.Consumer{},
		brokers:       brokers,
		topic:         topic,
	}
}

func (c *KafkaEventConsumer) RegisterEventHandler(eventHandler EventHandler) {
	c.eventHandlers = append(c.eventHandlers, eventHandler)
}

func (c *KafkaEventConsumer) Start(ctx context.Context) error {
	for _, handler := range c.eventHandlers {
		eventHandler := handler

		kafkaConsumer, err := kafka.NewConsumer(
			c.brokers,
			c.topic,
			handler.Name(),
		)
		if err != nil {
			return fmt.Errorf("failed to create kafka consumer: %w", err)
		}

		c.consumers = append(c.consumers, kafkaConsumer)

		go kafkaConsumer.Consume(ctx, func(record kafka.Record) {
			event, err := domain.UnmarshalEvent(record.Value)
			if err != nil {
				fmt.Println("failed to unmarshal event: ", err)
				return
			}

			if !slices.Contains(eventHandler.Events(), event.EventType) {
				return
			}

			if err := eventHandler.HandleEvent(ctx, event.Data); err != nil {
				fmt.Println("failed to handle event: ", err)
			}
		})
	}

	return nil
}

func (c *KafkaEventConsumer) Stop() {
	for _, consumer := range c.consumers {
		consumer.Close()
	}
}
