package kafka

import (
	"context"

	"github.com/cybre/home-inventory/internal/infrastructure"
	"github.com/cybre/home-inventory/inventory/app"
)

func NewKafkaTransport(ctx context.Context, eventMessaging *infrastructure.KafkaEventMessaging) error {
	if err := eventMessaging.ConsumeEvents(ctx, app.NewItemProjector()); err != nil {
		panic(err)
	}

	return nil
}
