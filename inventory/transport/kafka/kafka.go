package kafka

import (
	"context"

	"github.com/cybre/home-inventory/internal/infrastructure"
	"github.com/cybre/home-inventory/inventory/app/household"
)

func NewKafkaTransport(ctx context.Context, eventMessaging *infrastructure.KafkaEventMessaging) error {
	if err := eventMessaging.ConsumeEvents(ctx, household.NewItemProjector()); err != nil {
		panic(err)
	}

	return nil
}
