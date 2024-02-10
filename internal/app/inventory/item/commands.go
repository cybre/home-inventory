package item

import (
	"github.com/cybre/home-inventory/internal/app/shared"
	"github.com/cybre/home-inventory/pkg/domain"
)

type AddItemCommand struct {
	ItemID   string
	Name     string
	Barcode  string
	Quantity uint
}

func (c AddItemCommand) AggregateType() domain.AggregateType {
	return shared.InventoryAggregateType
}

func (c AddItemCommand) AggregateID() domain.AggregateID {
	return domain.AggregateID(c.ItemID)
}

type UpdateItemCommand struct {
	ItemID   string
	Name     string
	Barcode  string
	Quantity uint
}

func (c UpdateItemCommand) AggregateType() domain.AggregateType {
	return shared.InventoryAggregateType
}

func (c UpdateItemCommand) AggregateID() domain.AggregateID {
	return domain.AggregateID(c.ItemID)
}
