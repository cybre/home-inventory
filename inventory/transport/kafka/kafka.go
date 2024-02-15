package kafka

import (
	"context"

	"github.com/cybre/home-inventory/internal/infrastructure"
	"github.com/cybre/home-inventory/inventory/app/household"
	"github.com/cybre/home-inventory/inventory/app/user"
)

func NewKafkaTransport(ctx context.Context, eventMessaging *infrastructure.KafkaEventMessaging, loginInfoRepository *user.LoginInfoRepository, oneTimeLoginRepository *user.OneTimeLoginRepository) error {
	if err := eventMessaging.ConsumeEvents(ctx, user.NewLoginInfoProjector(loginInfoRepository)); err != nil {
		panic(err)
	}

	if err := eventMessaging.ConsumeEvents(ctx, user.NewOneTimeLoginTokenHandler(oneTimeLoginRepository)); err != nil {
		panic(err)
	}

	if err := eventMessaging.ConsumeEvents(ctx, household.NewItemProjector()); err != nil {
		panic(err)
	}

	return nil
}
