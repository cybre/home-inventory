package household

import (
	"github.com/cybre/home-inventory/pkg/domain"
)

type CreateHouseholdCommand struct {
	HouseholdID string
	UserID      string
	Name        string
}

func (c CreateHouseholdCommand) AggregateType() domain.AggregateType {
	return HouseholdAggregateType
}

func (c CreateHouseholdCommand) AggregateID() domain.AggregateID {
	return domain.AggregateID(c.HouseholdID)
}

type AddRoomCommand struct {
	HouseholdID string
	RoomID      string
	Name        string
}

func (c AddRoomCommand) AggregateType() domain.AggregateType {
	return HouseholdAggregateType
}

func (c AddRoomCommand) AggregateID() domain.AggregateID {
	return domain.AggregateID(c.HouseholdID)
}

type AddItemCommand struct {
	HouseholdID string
	RoomID      string
	ItemID      string
	Name        string
	Barcode     string
	Quantity    uint
}

func (c AddItemCommand) AggregateType() domain.AggregateType {
	return HouseholdAggregateType
}

func (c AddItemCommand) AggregateID() domain.AggregateID {
	return domain.AggregateID(c.HouseholdID)
}

type UpdateItemCommand struct {
	HouseholdID string
	RoomID      string
	ItemID      string
	Name        string
	Barcode     string
	Quantity    uint
}

func (c UpdateItemCommand) AggregateType() domain.AggregateType {
	return HouseholdAggregateType
}

func (c UpdateItemCommand) AggregateID() domain.AggregateID {
	return domain.AggregateID(c.HouseholdID)
}
