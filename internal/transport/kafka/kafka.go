package kafka

import (
	"context"

	"github.com/cybre/home-inventory/internal/app/inventory/item"
	"github.com/cybre/home-inventory/internal/infrastructure"
)

func NewKafkaTransport(ctx context.Context, eventMessaging *infrastructure.KafkaEventMessaging) error {
	if err := eventMessaging.ConsumeEvents(ctx, item.NewItemProjector()); err != nil {
		panic(err)
	}

	return nil
}
