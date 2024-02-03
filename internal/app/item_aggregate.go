package app

import (
	"context"
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
	domain.AggregateContext

	name string
}

func NewItemAggregate(aggregateContext domain.AggregateContext) domain.AggregateRoot {
	return &ItemAggregate{
		AggregateContext: aggregateContext,
	}
}

func (a *ItemAggregate) ApplyEvent(event domain.EventData) {
	switch e := event.(type) {
	case ItemAddedEvent:
		a.applyItemAddedEvent(e)
	case ItemUpdatedEvent:
		a.applyItemUpdatedEvent(e)
	default:
		panic("unknown event type")
	}
}

func (a *ItemAggregate) HandleCommand(ctx context.Context, command domain.Command) error {
	switch c := command.(type) {
	case AddItemCommand:
		return a.handleAddItemCommand(ctx, c)
	case UpdateItemCommand:
		return a.handleUpdateItemCommand(ctx, c)
	default:
		return domain.ErrUnknownCommand
	}
}

func (a *ItemAggregate) handleAddItemCommand(ctx context.Context, c AddItemCommand) error {
	if a.Version() != initialAggregateVersion {
		return ErrItemAlreadyExists
	}

	a.StoreEvent(ItemAddedEvent{
		Name: c.Name,
	})

	return nil
}

func (a *ItemAggregate) handleUpdateItemCommand(ctx context.Context, c UpdateItemCommand) error {
	if a.Version() == initialAggregateVersion {
		return ErrItemDoesNotExist
	}

	a.StoreEvent(ItemUpdatedEvent{
		Name: c.Name,
	})

	return nil
}

func (a *ItemAggregate) applyItemAddedEvent(e ItemAddedEvent) {
	a.name = e.Name
}

func (a *ItemAggregate) applyItemUpdatedEvent(e ItemUpdatedEvent) {
	a.name = e.Name
}
