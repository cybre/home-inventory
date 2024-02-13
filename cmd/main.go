package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/cybre/home-inventory/internal/infrastructure"
	"github.com/cybre/home-inventory/inventory/app"
	"github.com/cybre/home-inventory/inventory/domain/household"
	httptransport "github.com/cybre/home-inventory/inventory/transport/http"
	kafkatransport "github.com/cybre/home-inventory/inventory/transport/kafka"
	"github.com/cybre/home-inventory/pkg/cassandra"
	es "github.com/cybre/home-inventory/pkg/eventsourcing"
)

var kafkaBrokers = []string{"127.0.0.1:9092"}

const eventsTopic = "home-inventory.events"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	es.RegisterAggregateRoot(household.HouseholdAggregateType, household.NewHouseholdAggregate)
	es.RegisterEvent(household.ItemAddedEvent{})
	es.RegisterEvent(household.ItemUpdatedEvent{})

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

	if err := kafkatransport.NewKafkaTransport(ctx, eventMessaging); err != nil {
		panic(err)
	}

	if err := httptransport.NewHTTPTransport(ctx, householdService); err != nil {
		panic(err)
	}
}
