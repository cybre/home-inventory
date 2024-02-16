package kafka

import (
	"context"

	"github.com/cybre/home-inventory/internal/infrastructure"
	"github.com/cybre/home-inventory/inventory/app/household"
	"github.com/cybre/home-inventory/inventory/app/user"
)

func NewKafkaTransport(ctx context.Context, eventMessaging *infrastructure.KafkaEventMessaging, loginInfoRepository *user.LoginInfoRepository, loginTokenRepository *user.LoginTokenRepository) error {
	if err := eventMessaging.ConsumeEvents(ctx, user.NewLoginInfoProjector(loginInfoRepository)); err != nil {
		panic(err)
	}

	if err := eventMessaging.ConsumeEvents(ctx, user.NewLoginTokenHandler(loginTokenRepository)); err != nil {
		panic(err)
	}

	if err := eventMessaging.ConsumeEvents(ctx, household.NewItemProjector()); err != nil {
		panic(err)
	}

	return nil
}
