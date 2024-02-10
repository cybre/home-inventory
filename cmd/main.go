package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cybre/home-inventory/internal/app/item"
	"github.com/cybre/home-inventory/internal/infrastructure"
	httptransport "github.com/cybre/home-inventory/internal/transport/http"
	"github.com/cybre/home-inventory/pkg/cassandra"
	"github.com/cybre/home-inventory/pkg/domain"
	"github.com/cybre/home-inventory/pkg/kafka"
	"github.com/google/uuid"
)

var kafkaBrokers = []string{"127.0.0.1:9092"}

const eventsTopic = "home-inventory.events"

func main() {
	domain.RegisterAggregateRoot(item.ItemAggregateType, item.NewItemAggregate)
	domain.RegisterEvent(item.ItemAddedEvent{})
	domain.RegisterEvent(item.ItemUpdatedEvent{})

	cassandraSession, err := cassandra.NewSession([]string{"127.0.0.1:9042"}, "home_inventory")
	if err != nil {
		panic(err)
	}
	defer cassandraSession.Close()

	kafkaProducer, err := kafka.NewProducer(kafkaBrokers, eventsTopic)
	if err != nil {
		panic(err)
	}
	defer kafkaProducer.Close()

	eventPublisher := infrastructure.NewKafkaEventPublisher(kafkaProducer)
	eventStore := infrastructure.NewCassandraEventStore(cassandraSession)
	commandBus := domain.NewCommandBus(eventStore, eventPublisher)

	eventConsumer := infrastructure.NewKafkaEventConsumer(kafkaBrokers, eventsTopic)
	eventConsumer.RegisterEventHandler(item.NewItemProjector())
	if err := eventConsumer.Start(context.Background()); err != nil {
		panic(err)
	}
	defer eventConsumer.Stop()

	itemService := item.NewItemService(commandBus)

	itemId := uuid.NewString()

	if err := itemService.AddItem(context.Background(), item.AddItemCommandData{
		ItemID: itemId,
		Name:   "Test Item",
	}); err != nil {
		panic(err)
	}

	time.AfterFunc(30*time.Second, func() {
		if err := itemService.UpdateItem(context.Background(), item.UpdateItemCommandData{
			ItemID: itemId,
			Name:   "Test Item Updated",
		}); err != nil {
			fmt.Println(err)
		}
	})

	if err := httptransport.NewHTTPTransport(itemService); err != nil {
		panic(err)
	}
}
