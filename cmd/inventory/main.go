package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/cybre/home-inventory/internal/infrastructure"
	apphousehold "github.com/cybre/home-inventory/services/inventory/app/household"
	"github.com/cybre/home-inventory/services/inventory/domain/household"

	"github.com/cybre/home-inventory/internal/cassandra"
	es "github.com/cybre/home-inventory/internal/eventsourcing"
	"github.com/cybre/home-inventory/internal/logging"
	httptransport "github.com/cybre/home-inventory/services/inventory/transport/http"
	kafkatransport "github.com/cybre/home-inventory/services/inventory/transport/kafka"
)

var (
	kafkaBrokers   = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	cassandraHosts = strings.Split(os.Getenv("CASSANDRA_HOSTS"), ",")
	serverAddress  = os.Getenv("SERVER_ADDRESS")
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

	eventMessaging, err := infrastructure.NewKafkaEventMessaging(kafkaBrokers, eventsTopic, logger)
	if err != nil {
		panic(err)
	}
	defer eventMessaging.Close()

	es.RegisterAggregateRoot(household.HouseholdAggregateType, household.NewHouseholdAggregate)
	es.RegisterEvent(household.HouseholdCreatedEvent{})
	es.RegisterEvent(household.HouseholdUpdatedEvent{})
	es.RegisterEvent(household.HouseholdDeletedEvent{})
	es.RegisterEvent(household.RoomAddedEvent{})
	es.RegisterEvent(household.RoomUpdatedEvent{})
	es.RegisterEvent(household.RoomDeletedEvent{})

	eventStore, err := infrastructure.NewCassandraEventStore(cassandraSession)
	if err != nil {
		panic(err)
	}
	commandBus := es.NewCommandBus(eventStore, eventMessaging)

	userHouseholdRepository := apphousehold.NewUserHouseholdRepository(cassandraSession)
	householdService := apphousehold.NewHouseholdService(commandBus, userHouseholdRepository)

	if err := kafkatransport.NewKafkaTransport(ctx, eventMessaging, userHouseholdRepository); err != nil {
		panic(err)
	}

	if err := httptransport.NewHTTPTransport(ctx, serverAddress, householdService); err != nil {
		panic(err)
	}
}
