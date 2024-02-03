package app

import (
	"fmt"

	"github.com/cybre/home-inventory/pkg/domain"
)

const (
	initialAggregateVersion                      = 0
	ItemAggregateType       domain.AggregateType = "ItemAggregate"
)

var (
	ErrItemAlreadyExists = fmt.Errorf("item already exists")
	ErrItemDoesNotExist  = fmt.Errorf("item does not exist")
)

type ItemAggregate struct {
	*domain.AggregateContext

	name string
}

func NewItemAggregate(aggregateContext *domain.AggregateContext) domain.AggregateRoot {
	return &ItemAggregate{
		AggregateContext: aggregateContext,
	}
}

func (a *ItemAggregate) ApplyEvent(event domain.EventData) {
	switch e := event.(type) {
	case ItemAddedEvent:
		a.name = e.Name
	case ItemUpdatedEvent:
		a.name = e.Name
	default:
		panic("unknown event type")
	}
}

func (a *ItemAggregate) HandleCommand(command domain.Command) error {
	switch c := command.(type) {
	case AddItemCommand:
		if a.Version() != initialAggregateVersion {
			return ErrItemAlreadyExists
		}

		a.StoreEvent(ItemAddedEvent{
			Name: c.Name,
		})

		return nil
	case UpdateItemCommand:
		if a.Version() == initialAggregateVersion {
			return ErrItemDoesNotExist
		}

		a.StoreEvent(ItemUpdatedEvent{
			Name: c.Name,
		})

		return nil
	default:
		panic("unknown command type")
	}
}
