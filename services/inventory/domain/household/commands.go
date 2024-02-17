package household

import es "github.com/cybre/home-inventory/internal/eventsourcing"

type CreateHouseholdCommand struct {
	HouseholdID string
	UserID      string
	Name        string
}

func (c CreateHouseholdCommand) AggregateType() es.AggregateType {
	return HouseholdAggregateType
}

func (c CreateHouseholdCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.HouseholdID)
}

type AddRoomCommand struct {
	HouseholdID string
	RoomID      string
	Name        string
}

func (c AddRoomCommand) AggregateType() es.AggregateType {
	return HouseholdAggregateType
}

func (c AddRoomCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.HouseholdID)
}

type AddItemCommand struct {
	HouseholdID string
	RoomID      string
	ItemID      string
	Name        string
	Barcode     string
	Quantity    uint
}

func (c AddItemCommand) AggregateType() es.AggregateType {
	return HouseholdAggregateType
}

func (c AddItemCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.HouseholdID)
}

type UpdateItemCommand struct {
	HouseholdID string
	RoomID      string
	ItemID      string
	Name        string
	Barcode     string
	Quantity    uint
}

func (c UpdateItemCommand) AggregateType() es.AggregateType {
	return HouseholdAggregateType
}

func (c UpdateItemCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.HouseholdID)
}
