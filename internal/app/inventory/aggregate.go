package inventory

import (
	"context"
	"fmt"

	"github.com/cybre/home-inventory/internal/app/inventory/household"
	"github.com/cybre/home-inventory/internal/app/inventory/item"
	s "github.com/cybre/home-inventory/internal/app/shared"
	"github.com/cybre/home-inventory/pkg/domain"
)

const (
	initialAggregateVersion = 0
)

var (
	ErrItemAlreadyExists = fmt.Errorf("item already exists")
	ErrItemDoesNotExist  = fmt.Errorf("item does not exist")
)

type InventoryAggregate struct {
	domain.AggregateContext

	Household household.Household
	Items     item.Items
}

func NewInventoryAggregate(aggregateContext domain.AggregateContext) domain.AggregateRoot {
	return &InventoryAggregate{
		AggregateContext: aggregateContext,
	}
}

func (a *InventoryAggregate) ApplyEvent(event domain.EventData) {
	switch e := event.(type) {
	case item.ItemAddedEvent:
		a.applyItemAddedEvent(e)
	case item.ItemUpdatedEvent:
		a.applyItemUpdatedEvent(e)
	default:
		panic("unknown event type")
	}
}

func (a *InventoryAggregate) HandleCommand(ctx context.Context, command domain.Command) ([]domain.EventData, error) {
	switch c := command.(type) {
	case item.AddItemCommand:
		return a.handleAddItemCommand(ctx, c)
	case item.UpdateItemCommand:
		return a.handleUpdateItemCommand(ctx, c)
	default:
		return nil, domain.ErrUnknownCommand
	}
}

func (a *InventoryAggregate) handleAddItemCommand(ctx context.Context, c item.AddItemCommand) ([]domain.EventData, error) {
	itemID, err := item.NewItemID(c.ItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create item id: %w", err)
	}

	if a.Items.Exists(itemID) {
		return nil, ErrItemAlreadyExists
	}

	item, err := item.NewItem(c.ItemID, c.Name, c.Barcode, c.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return s.Events(item.ToItemAddedEvent())
}

func (a *InventoryAggregate) handleUpdateItemCommand(ctx context.Context, c item.UpdateItemCommand) ([]domain.EventData, error) {
	itemID, err := item.NewItemID(c.ItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to create item id: %w", err)
	}

	item, ok := a.Items.Get(itemID)
	if !ok {
		return nil, ErrItemDoesNotExist
	}

	item, err = item.Update(c.Name, c.Barcode, c.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return s.Events(item.ToItemUpdatedEvent())
}

func (a *InventoryAggregate) applyItemAddedEvent(e item.ItemAddedEvent) {
	item, _ := item.NewItemFromEvent(e)

	a.Items.Add(item)
}

func (a *InventoryAggregate) applyItemUpdatedEvent(e item.ItemUpdatedEvent) {
	item, ok := a.Items.Get(item.ItemID(e.ID))
	if ok {
		item, _ = item.UpdateFromEvent(e)
		a.Items.Update(item)
	}
}
