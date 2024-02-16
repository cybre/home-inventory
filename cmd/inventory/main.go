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
	loginTokenRepository := appuser.NewLoginTokenRepository(cassandraSession)

	emailUniquenessService := userservices.NewEmailUniquenessService(loginInfoRepository)

	es.RegisterAggregateRoot(household.HouseholdAggregateType, household.NewHouseholdAggregate)
	es.RegisterEvent(household.HouseholdCreatedEvent{})
	es.RegisterEvent(household.RoomAddedEvent{})
	es.RegisterEvent(household.ItemAddedEvent{})
	es.RegisterEvent(household.ItemUpdatedEvent{})

	es.RegisterAggregateRoot(user.UserAggregateType, user.NewUserAggregate(emailUniquenessService))
	es.RegisterEvent(user.UserCreatedEvent{})
	es.RegisterEvent(user.UserLoginTokenGeneratedEvent{})
	es.RegisterEvent(user.UserLoginEvent{})

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
	userService := appuser.NewUserService(commandBus, loginInfoRepository, loginTokenRepository)

	if err := kafkatransport.NewKafkaTransport(ctx, eventMessaging, loginInfoRepository, loginTokenRepository); err != nil {
		panic(err)
	}

	if err := httptransport.NewHTTPTransport(ctx, householdService, userService); err != nil {
		panic(err)
	}
}
