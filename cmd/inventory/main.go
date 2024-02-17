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

	es.RegisterAggregateRoot(household.HouseholdAggregateType, household.NewHouseholdAggregate)
	es.RegisterEvent(household.HouseholdCreatedEvent{})
	es.RegisterEvent(household.RoomAddedEvent{})
	es.RegisterEvent(household.ItemAddedEvent{})
	es.RegisterEvent(household.ItemUpdatedEvent{})

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

	if err := kafkatransport.NewKafkaTransport(ctx, eventMessaging); err != nil {
		panic(err)
	}

	if err := httptransport.NewHTTPTransport(ctx, householdService); err != nil {
		panic(err)
	}
}
