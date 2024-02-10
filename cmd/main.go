package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cybre/home-inventory/internal/app/item"
	"github.com/cybre/home-inventory/internal/infrastructure"
	"github.com/cybre/home-inventory/internal/infrastructure/cassandra"
	"github.com/cybre/home-inventory/internal/infrastructure/kafka"
	httptransport "github.com/cybre/home-inventory/internal/transport/http"
	"github.com/cybre/home-inventory/pkg/domain"
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
	eventConsumer.RegisterEventHandler(item.NewItemProjector())
	if err := eventConsumer.Start(context.Background()); err != nil {
		panic(err)
	}
	defer eventConsumer.Stop()

	itemService := item.NewItemService(commandBus)

	// if err := itemService.AddItem(context.Background(), app.AddItemCommandData{
	// 	ItemID: "65396437-3930-3039-2d35-6132352d3433",
	// 	Name:   "Test Item",
	// }); err != nil {
	// 	panic(err)
	// }

	if err := itemService.UpdateItem(context.Background(), item.UpdateItemCommandData{
		ItemID: "65396437-3930-3039-2d35-6132352d3433",
		Name:   "Test Item 3",
	}); err != nil {
		panic(err)
	}

	time.AfterFunc(30*time.Second, func() {
		if err := itemService.UpdateItem(context.Background(), item.UpdateItemCommandData{
			ItemID: "65396437-3930-3039-2d35-6132352d3433",
			Name:   "Test Item 4",
		}); err != nil {
			fmt.Println(err)
		}
	})

	httpTransport := httptransport.NewHTTPTransport(itemService)
	http.ListenAndServe(":3000", httpTransport)
}
