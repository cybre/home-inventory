package kafka

import (
	"context"

	"github.com/cybre/home-inventory/internal/infrastructure"
	"github.com/cybre/home-inventory/services/inventory/app/household"
)

func NewKafkaTransport(ctx context.Context, eventMessaging *infrastructure.KafkaEventMessaging, userHouseholdRepository *household.UserHouseholdRepository) error {
	if err := eventMessaging.ConsumeEvents(ctx, household.NewUserHouseholdProjector(userHouseholdRepository)); err != nil {
		panic(err)
	}

	if err := eventMessaging.ConsumeEvents(ctx, household.NewItemProjector()); err != nil {
		panic(err)
	}

	return nil
}
