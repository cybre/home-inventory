package kafka

import (
	"context"

	"github.com/cybre/home-inventory/internal/app/item"
	"github.com/cybre/home-inventory/internal/infrastructure"
)

type EventConsumer interface {
	ConsumeEvents(context.Context, infrastructure.EventHandler) error
}

func NewKafkaTransport(ctx context.Context, consumer EventConsumer) error {
	if err := consumer.ConsumeEvents(ctx, item.NewItemProjector()); err != nil {
		panic(err)
	}

	return nil
}
