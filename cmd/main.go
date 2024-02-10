package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/cybre/home-inventory/internal/app/item"
	"github.com/cybre/home-inventory/internal/infrastructure"
	httptransport "github.com/cybre/home-inventory/internal/transport/http"
	kafkatransport "github.com/cybre/home-inventory/internal/transport/kafka"
	"github.com/cybre/home-inventory/pkg/cassandra"
	"github.com/cybre/home-inventory/pkg/domain"
	"github.com/cybre/home-inventory/pkg/kafka"
	"github.com/google/uuid"
)

var kafkaBrokers = []string{"127.0.0.1:9092"}

const eventsTopic = "home-inventory.events"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	domain.RegisterAggregateRoot(item.ItemAggregateType, item.NewItemAggregate)
	domain.RegisterEvent(item.ItemAddedEvent{})
	domain.RegisterEvent(item.ItemUpdatedEvent{})

	cassandraSession, err := cassandra.NewSession([]string{"127.0.0.1:9042"}, "home_inventory")
	if err != nil {
		panic(err)
	}
	defer cassandraSession.Close()

	kafkaProducer, err := kafka.NewProducer(kafkaBrokers, eventsTopic, logger)
	if err != nil {
		panic(err)
	}
	defer kafkaProducer.Close()

	eventMessaging, err := infrastructure.NewKafkaEventMessaging(kafkaBrokers, eventsTopic, logger)
	if err != nil {
		panic(err)
	}
	defer eventMessaging.Close()

	eventStore := infrastructure.NewCassandraEventStore(cassandraSession)
	commandBus := domain.NewCommandBus(eventStore, eventMessaging)

	itemService := item.NewItemService(commandBus)

	itemId := uuid.NewString()
	if err := itemService.AddItem(ctx, item.AddItemCommandData{
		ItemID: itemId,
		Name:   "Test Item",
	}); err != nil {
		panic(err)
	}
	time.AfterFunc(30*time.Second, func() {
		if err := itemService.UpdateItem(ctx, item.UpdateItemCommandData{
			ItemID: itemId,
			Name:   "Test Item Updated",
		}); err != nil {
			slog.Error("failed to update item", slog.Any("error", err))
		}
	})

	if err := kafkatransport.NewKafkaTransport(ctx, eventMessaging); err != nil {
		panic(err)
	}

	if err := httptransport.NewHTTPTransport(ctx, itemService); err != nil {
		panic(err)
	}
}
