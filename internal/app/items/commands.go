package items

import "github.com/cybre/home-inventory/pkg/domain"

type AddItemCommand struct {
	ItemID string
	Name   string
}

func (c AddItemCommand) AggregateType() domain.AggregateType {
	return ItemAggregateType
}

func (c AddItemCommand) AggregateID() domain.AggregateID {
	return domain.AggregateID(c.ItemID)
}

type UpdateItemCommand struct {
	ItemID string
	Name   string
}

func (c UpdateItemCommand) AggregateType() domain.AggregateType {
	return ItemAggregateType
}

func (c UpdateItemCommand) AggregateID() domain.AggregateID {
	return domain.AggregateID(c.ItemID)
}
