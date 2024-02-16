package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/cybre/home-inventory/internal/infrastructure"
	apphousehold "github.com/cybre/home-inventory/inventory/app/household"
	appuser "github.com/cybre/home-inventory/inventory/app/user"
	"github.com/cybre/home-inventory/inventory/domain/household"
	"github.com/cybre/home-inventory/inventory/domain/user"
	userservices "github.com/cybre/home-inventory/inventory/domain/user/services"

	httptransport "github.com/cybre/home-inventory/inventory/transport/http"
	kafkatransport "github.com/cybre/home-inventory/inventory/transport/kafka"
	"github.com/cybre/home-inventory/pkg/cassandra"
	es "github.com/cybre/home-inventory/pkg/eventsourcing"
	"github.com/cybre/home-inventory/pkg/logging"
)

var (
	kafkaBrokers   = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	cassandraHosts = strings.Split(os.Getenv("CASSANDRA_HOSTS"), ",")
)

const (
	eventsTopic = "inventory.events"
	serviceName = "inventory"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("service", serviceName)
	slog.SetDefault(logger)

	ctx = logging.WithLogger(ctx, logger)

	cassandraSession, err := cassandra.NewSession(cassandraHosts, serviceName)
	if err != nil {
		panic(err)
	}
	defer cassandraSession.Close()

	loginInfoRepository := appuser.NewLoginInfoRepository(cassandraSession)
	oneTimeLoginRepository := appuser.NewOneTimeLoginRepository(cassandraSession)

	emailUniquenessService := userservices.NewEmailUniquenessService(loginInfoRepository)

	es.RegisterAggregateRoot(household.HouseholdAggregateType, household.NewHouseholdAggregate)
	es.RegisterEvent(household.HouseholdCreatedEvent{})
	es.RegisterEvent(household.RoomAddedEvent{})
	es.RegisterEvent(household.ItemAddedEvent{})
	es.RegisterEvent(household.ItemUpdatedEvent{})

	es.RegisterAggregateRoot(user.UserAggregateType, user.NewUserAggregate(emailUniquenessService))
	es.RegisterEvent(user.UserCreatedEvent{})
	es.RegisterEvent(user.UserOneTimeTokenGeneratedEvent{})

	eventMessaging, err := infrastructure.NewKafkaEventMessaging(kafkaBrokers, eventsTopic, logger)
	if err != nil {
		panic(err)
	}
	defer eventMessaging.Close()

	eventStore, err := infrastructure.NewCassandraEventStore(cassandraSession)
	if err != nil {
		panic(err)
	}
	commandBus := es.NewCommandBus(eventStore, eventMessaging)

	householdService := apphousehold.NewHouseholdService(commandBus)
	userService := appuser.NewUserService(commandBus, loginInfoRepository)

	if err := kafkatransport.NewKafkaTransport(ctx, eventMessaging, loginInfoRepository, oneTimeLoginRepository); err != nil {
		panic(err)
	}

	// time.Sleep(5 * time.Second)

	// userId := uuid.NewString()

	// if err := userService.CreateUser(ctx, shared.CreateUserCommandData{
	// 	UserID:    userId,
	// 	FirstName: "Stefan",
	// 	LastName:  "Ric",
	// 	Email:     "stfric369@gmail.com",
	// }); err != nil {
	// 	panic(err)
	// }

	// if err := userService.GenerateOneTimeToken(ctx, shared.GenerateOneTimeTokenCommandData{
	// 	Email: "stfric369@gmail.com",
	// }); err != nil {
	// 	panic(err)
	// }

	// householdId := uuid.NewString()

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

	// itemId := uuid.NewString()

	// if err := householdService.AddItem(ctx, shared.AddItemCommandData{
	// 	HouseholdID: householdId,
	// 	RoomID:      roomId,
	// 	ItemID:      itemId,
	// 	Name:        "Test Item 123",
	// 	Barcode:     "1234567890",
	// 	Quantity:    1,
	// }); err != nil {
	// 	panic(err)
	// }

	// if err := householdService.UpdateItem(ctx, shared.UpdateItemCommandData{
	// 	HouseholdID: householdId,
	// 	RoomID:      roomId,
	// 	ItemID:      itemId,
	// 	Name:        "Test Item 1234",
	// 	Barcode:     "1234567890",
	// 	Quantity:    1,
	// }); err != nil {
	// 	panic(err)
	// }

	if err := httptransport.NewHTTPTransport(ctx, householdService, userService); err != nil {
		panic(err)
	}
}
