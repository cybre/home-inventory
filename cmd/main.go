package main

import (
	"context"
	"net/http"

	"github.com/cybre/home-inventory/internal/app"
	"github.com/cybre/home-inventory/internal/infrastructure"
	httptransport "github.com/cybre/home-inventory/internal/transport/http"
	"github.com/cybre/home-inventory/pkg/domain"
)

func main() {
	cassandraSession, err := infrastructure.NewCassandraSession([]string{"127.0.0.1:9042"}, "home_inventory")
	if err != nil {
		panic(err)
	}
	defer cassandraSession.Close()

	eventStore := infrastructure.NewCassandraEventStore(cassandraSession)
	commandBus := domain.NewCommandBus(eventStore)

	domain.RegisterAggregateRoot(app.ItemAggregateType, app.NewItemAggregate)
	domain.RegisterEvent(app.ItemAddedEvent{})
	domain.RegisterEvent(app.ItemUpdatedEvent{})

	itemService := app.NewItemService(commandBus)

	// if err := itemService.AddItem(context.Background(), app.AddItemCommandData{
	// 	ItemID: "65396437-3930-3039-2d35-6132352d3433",
	// 	Name:   "Test Item",
	// }); err != nil {
	// 	panic(err)
	// }

	if err := itemService.UpdateItem(context.Background(), app.UpdateItemCommandData{
		ItemID: "65396437-3930-3039-2d35-6132352d3433",
		Name:   "Test Item 3",
	}); err != nil {
		panic(err)
	}

	httpTransport := httptransport.NewHTTPTransport(itemService)
	http.ListenAndServe(":3000", httpTransport)
}
