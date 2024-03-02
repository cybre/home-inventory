package household

import es "github.com/cybre/home-inventory/internal/eventsourcing"

type CreateHouseholdCommand struct {
	HouseholdID string
	UserID      string
	Name        string
	Location    string
	Description string
}

func (c CreateHouseholdCommand) AggregateType() es.AggregateType {
	return HouseholdAggregateType
}

func (c CreateHouseholdCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.HouseholdID)
}

type UpdateHouseholdCommand struct {
	HouseholdID string
	UserID      string
	Name        string
	Location    string
	Description string
}

func (c UpdateHouseholdCommand) AggregateType() es.AggregateType {
	return HouseholdAggregateType
}

func (c UpdateHouseholdCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.HouseholdID)
}

type DeleteHouseholdCommand struct {
	HouseholdID string
	UserID      string
}

func (c DeleteHouseholdCommand) AggregateType() es.AggregateType {
	return HouseholdAggregateType
}

func (c DeleteHouseholdCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.HouseholdID)
}

type AddRoomCommand struct {
	HouseholdID string
	UserID      string
	RoomID      string
	Name        string
}

func (c AddRoomCommand) AggregateType() es.AggregateType {
	return HouseholdAggregateType
}

func (c AddRoomCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.HouseholdID)
}

type UpdateRoomCommand struct {
	HouseholdID string
	UserID      string
	RoomID      string
	Name        string
}

func (c UpdateRoomCommand) AggregateType() es.AggregateType {
	return HouseholdAggregateType
}

func (c UpdateRoomCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.HouseholdID)
}

type DeleteRoomCommand struct {
	HouseholdID string
	UserID      string
	RoomID      string
}

func (c DeleteRoomCommand) AggregateType() es.AggregateType {
	return HouseholdAggregateType
}

func (c DeleteRoomCommand) AggregateID() es.AggregateID {
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
