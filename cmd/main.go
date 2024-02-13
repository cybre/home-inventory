package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/cybre/home-inventory/internal/infrastructure"
	"github.com/cybre/home-inventory/inventory/app"
	"github.com/cybre/home-inventory/inventory/domain/household"
	"github.com/cybre/home-inventory/inventory/domain/user"

	// "github.com/cybre/home-inventory/inventory/shared"
	httptransport "github.com/cybre/home-inventory/inventory/transport/http"
	kafkatransport "github.com/cybre/home-inventory/inventory/transport/kafka"
	"github.com/cybre/home-inventory/pkg/cassandra"
	es "github.com/cybre/home-inventory/pkg/eventsourcing"
	// "github.com/google/uuid"
)

var kafkaBrokers = []string{"127.0.0.1:9092"}

const eventsTopic = "home-inventory.events"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	es.RegisterAggregateRoot(household.HouseholdAggregateType, household.NewHouseholdAggregate)
	es.RegisterEvent(household.HouseholdCreatedEvent{})
	es.RegisterEvent(household.RoomAddedEvent{})
	es.RegisterEvent(household.ItemAddedEvent{})
	es.RegisterEvent(household.ItemUpdatedEvent{})
	es.RegisterEvent(user.UserCreatedEvent{})

	cassandraSession, err := cassandra.NewSession([]string{"127.0.0.1:9042"}, "home_inventory")
	if err != nil {
		panic(err)
	}
	defer cassandraSession.Close()

	eventMessaging, err := infrastructure.NewKafkaEventMessaging(kafkaBrokers, eventsTopic, logger)
	if err != nil {
		panic(err)
	}
	defer eventMessaging.Close()

	eventStore := infrastructure.NewCassandraEventStore(cassandraSession)
	commandBus := es.NewCommandBus(eventStore, eventMessaging)

	householdService := app.NewHouseholdService(commandBus)

	// householdId := uuid.NewString()
	// userId := uuid.NewString()

	// if err := householdService.CreateHousehold(ctx, shared.CreateHouseholdCommandData{
	// 	HouseholdID: householdId,
	// 	UserID:      userId,
	// 	Name:        "Test Household",
	// }); err != nil {
	// 	panic(err)
	// }

	// roomId := uuid.NewString()

	// if err := householdService.AddRoom(ctx, shared.AddRoomCommandData{
	// 	HouseholdID: householdId,
	// 	RoomID:      roomId,
	// 	Name:        "Test Room",
	// }); err != nil {
	// 	panic(err)
	// }

	// if err := householdService.AddItem(ctx, shared.AddItemCommandData{
	// 	HouseholdID: householdId,
	// 	RoomID:      roomId,
	// 	ItemID:      uuid.NewString(),
	// 	Name:        "Test Item",
	// 	Barcode:     "1234567890",
	// 	Quantity:    1,
	// }); err != nil {
	// 	panic(err)
	// }

	if err := kafkatransport.NewKafkaTransport(ctx, eventMessaging); err != nil {
		panic(err)
	}

	if err := httptransport.NewHTTPTransport(ctx, householdService); err != nil {
		panic(err)
	}
}
