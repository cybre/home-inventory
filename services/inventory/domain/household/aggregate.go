package household

import (
	"context"
	"fmt"
	"time"

	es "github.com/cybre/home-inventory/internal/eventsourcing"
	"github.com/cybre/home-inventory/services/inventory/domain/common"
	c "github.com/cybre/home-inventory/services/inventory/domain/common"
)

const (
	HouseholdAggregateType  es.AggregateType = "HouseholdAggregate"
	initialAggregateVersion                  = 0
)

type HouseholdService interface {
	GetHouseholdCount(ctx context.Context, userID common.UserID) (uint, error)
	CheckHouseholdNameAvailability(ctx context.Context, userID common.UserID, name HouseholdName) (bool, error)
}

type HouseholdAgregate struct {
	es.AggregateContext

	householdService HouseholdService

	UserID      c.UserID
	Name        HouseholdName
	Location    HouseholdLocation
	Description HouseholdDescription
	Order       HouseholdOrder

	Rooms Rooms
}

func NewHouseholdAggregate(householdService HouseholdService) es.AggregateRootFactoryFunc {
	return func(aggregateContext es.AggregateContext) es.AggregateRoot {
		return &HouseholdAgregate{
			AggregateContext: aggregateContext,
			householdService: householdService,
		}
	}
}

func (a *HouseholdAgregate) ApplyEvent(event es.EventData) {
	switch e := event.(type) {
	case HouseholdCreatedEvent:
		a.applyHouseholdCreatedEvent(e)
	case HouseholdUpdatedEvent:
		a.applyHouseholdUpdatedEvent(e)
	case RoomAddedEvent:
		a.applyRoomAddedEvent(e)
	case RoomUpdatedEvent:
		a.applyRoomUpdatedEvent(e)
	case ItemAddedEvent:
		a.applyItemAddedEvent(e)
	case ItemUpdatedEvent:
		a.applyItemUpdatedEvent(e)
	default:
		panic("unknown event type")
	}
}

func (a *HouseholdAgregate) HandleCommand(ctx context.Context, command es.Command) ([]es.EventData, error) {
	switch c := command.(type) {
	case CreateHouseholdCommand:
		return a.handleCreateHouseholdCommand(ctx, c)
	case UpdateHouseholdCommand:
		return a.handleUpdateHouseholdCommand(ctx, c)
	case AddRoomCommand:
		return a.handleAddRoomCommand(ctx, c)
	case UpdateRoomCommand:
		return a.handleUpdateRoomCommand(ctx, c)
	case AddItemCommand:
		return a.handleAddItemCommand(ctx, c)
	case UpdateItemCommand:
		return a.handleUpdateItemCommand(ctx, c)
	default:
		return nil, es.ErrUnknownCommand
	}
}

func (a *HouseholdAgregate) handleCreateHouseholdCommand(ctx context.Context, command CreateHouseholdCommand) ([]es.EventData, error) {
	if a.Version() != initialAggregateVersion {
		return nil, fmt.Errorf("household with provided ID already exists: %s", command.HouseholdID)
	}

	userId, err := c.NewUserID(command.UserID)
	if err != nil {
		return nil, err
	}

	name, err := NewHouseholdName(command.Name)
	if err != nil {
		return nil, err
	}

	nameAvailable, err := a.householdService.CheckHouseholdNameAvailability(ctx, userId, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check household name availability: %w", err)
	}
	if !nameAvailable {
		return nil, fmt.Errorf("household with name [%s] already exists", command.Name)
	}

	location, err := NewHouseholdLocation(command.Location)
	if err != nil {
		return nil, err
	}

	description, err := NewHouseholdDescription(command.Description)
	if err != nil {
		return nil, err
	}

	householdCount, err := a.householdService.GetHouseholdCount(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get household count: %w", err)
	}

	return c.Events(HouseholdCreatedEvent{
		HouseholdID: a.AggregateID().String(),
		UserID:      userId.String(),
		Name:        name.String(),
		Location:    location.String(),
		Description: description.String(),
		Order:       householdCount + 1,
		Timestamp:   time.Now().UnixMilli(),
	})
}

func (a *HouseholdAgregate) handleUpdateHouseholdCommand(ctx context.Context, command UpdateHouseholdCommand) ([]es.EventData, error) {
	if a.Version() == initialAggregateVersion {
		return nil, fmt.Errorf("household with provided ID does not exist: %s", command.HouseholdID)
	}

	name, err := NewHouseholdName(command.Name)
	if err != nil {
		return nil, err
	}

	location, err := NewHouseholdLocation(command.Location)
	if err != nil {
		return nil, err
	}

	description, err := NewHouseholdDescription(command.Description)
	if err != nil {
		return nil, err
	}

	return c.Events(HouseholdUpdatedEvent{
		HouseholdID: a.AggregateID().String(),
		UserID:      a.UserID.String(),
		Name:        name.String(),
		Location:    location.String(),
		Description: description.String(),
		Timestamp:   time.Now().UnixMilli(),
	})
}

func (a *HouseholdAgregate) handleAddRoomCommand(ctx context.Context, command AddRoomCommand) ([]es.EventData, error) {
	if a.Version() == initialAggregateVersion {
		return nil, fmt.Errorf("household with provided ID does not exist: %s", command.HouseholdID)
	}

	newRoom, err := NewRoom(command.RoomID, command.Name, uint(a.Rooms.Count()+1))
	if err != nil {
		return nil, err
	}

	if err := a.Rooms.Add(newRoom); err != nil {
		return nil, err
	}

	return c.Events(RoomAddedEvent{
		HouseholdID: a.AggregateID().String(),
		UserID:      a.UserID.String(),
		RoomID:      newRoom.ID.String(),
		Name:        newRoom.Name.String(),
		Order:       newRoom.Order.Uint(),
		Timestamp:   time.Now().UnixMilli(),
	})
}

func (a *HouseholdAgregate) handleUpdateRoomCommand(ctx context.Context, command UpdateRoomCommand) ([]es.EventData, error) {
	if a.Version() == initialAggregateVersion {
		return nil, fmt.Errorf("household with provided ID does not exist: %s", command.HouseholdID)
	}

	roomID, err := NewRoomID(command.RoomID)
	if err != nil {
		return nil, err
	}

	room, ok := a.Rooms.Get(roomID)
	if !ok {
		return nil, fmt.Errorf("room with ID %s does not exist", roomID)
	}

	room, err = room.Update(command.Name)
	if err != nil {
		return nil, err
	}

	return c.Events(RoomUpdatedEvent{
		HouseholdID: a.AggregateID().String(),
		UserID:      a.UserID.String(),
		RoomID:      room.ID.String(),
		Name:        room.Name.String(),
		Order:       room.Order.Uint(),
		Timestamp:   time.Now().UnixMilli(),
	})
}

func (a *HouseholdAgregate) handleAddItemCommand(ctx context.Context, command AddItemCommand) ([]es.EventData, error) {
	if a.Version() == initialAggregateVersion {
		return nil, fmt.Errorf("household with provided ID does not exist: %s", command.HouseholdID)
	}

	roomId, err := NewRoomID(command.RoomID)
	if err != nil {
		return nil, err
	}

	room, ok := a.Rooms.Get(roomId)
	if !ok {
		return nil, fmt.Errorf("room with ID %s does not exist", roomId)
	}

	item, err := NewItem(command.ItemID, command.Name, command.Barcode, command.Quantity)
	if err != nil {
		return nil, err
	}

	if err := room.Items.Add(item); err != nil {
		return nil, err
	}

	return c.Events(ItemAddedEvent{
		HouseholdID: a.AggregateID().String(),
		RoomID:      roomId.String(),
		ItemID:      item.ID.String(),
		Name:        item.Name.String(),
		Barcode:     item.Barcode.String(),
		Quantity:    item.Quantity.Uint(),
	})
}

func (a *HouseholdAgregate) handleUpdateItemCommand(ctx context.Context, command UpdateItemCommand) ([]es.EventData, error) {
	if a.Version() == initialAggregateVersion {
		return nil, fmt.Errorf("household with provided ID does not exist: %s", command.HouseholdID)
	}

	roomId, err := NewRoomID(command.RoomID)
	if err != nil {
		return nil, err
	}

	room, ok := a.Rooms.Get(roomId)
	if !ok {
		return nil, fmt.Errorf("room with ID %s does not exist", roomId)
	}

	itemId, err := NewItemID(command.ItemID)
	if err != nil {
		return nil, err
	}

	item, ok := room.Items.Get(itemId)
	if !ok {
		return nil, fmt.Errorf("item with ID %s does not exist", command.ItemID)
	}

	item, err = item.Update(command.Name, command.Barcode, command.Quantity)
	if err != nil {
		return nil, err
	}

	return c.Events(ItemUpdatedEvent{
		HouseholdID: a.AggregateID().String(),
		RoomID:      roomId.String(),
		ItemID:      item.ID.String(),
		Name:        item.Name.String(),
		Barcode:     item.Barcode.String(),
		Quantity:    item.Quantity.Uint(),
	})
}

func (a *HouseholdAgregate) applyHouseholdCreatedEvent(event HouseholdCreatedEvent) {
	a.UserID, _ = c.NewUserID(event.UserID)
	a.Name, _ = NewHouseholdName(event.Name)
	a.Location, _ = NewHouseholdLocation(event.Location)
	a.Description, _ = NewHouseholdDescription(event.Description)
	a.Order, _ = NewHouseholdOrder(event.Order)
	a.Rooms = NewRooms()
}

func (a *HouseholdAgregate) applyHouseholdUpdatedEvent(event HouseholdUpdatedEvent) {
	a.Name, _ = NewHouseholdName(event.Name)
	a.Location, _ = NewHouseholdLocation(event.Location)
	a.Description, _ = NewHouseholdDescription(event.Description)
}

func (a *HouseholdAgregate) applyRoomAddedEvent(event RoomAddedEvent) {
	newRoom, _ := NewRoom(event.RoomID, event.Name, event.Order)
	a.Rooms.Add(newRoom)
}

func (a *HouseholdAgregate) applyRoomUpdatedEvent(event RoomUpdatedEvent) {
	roomID, _ := NewRoomID(event.RoomID)
	room, _ := a.Rooms.Get(roomID)

	room, _ = room.Update(event.Name)

	a.Rooms.Update(room)
}

func (a *HouseholdAgregate) applyItemAddedEvent(event ItemAddedEvent) {
	roomID, _ := NewRoomID(event.RoomID)
	room, _ := a.Rooms.Get(roomID)

	item, _ := NewItem(event.ItemID, event.Name, event.Barcode, event.Quantity)
	room.Items.Add(item)
}

func (a *HouseholdAgregate) applyItemUpdatedEvent(event ItemUpdatedEvent) {
	roomID, _ := NewRoomID(event.RoomID)
	room, _ := a.Rooms.Get(roomID)

	itemID, _ := NewItemID(event.ItemID)
	item, _ := room.Items.Get(itemID)

	item, _ = item.Update(event.Name, event.Barcode, event.Quantity)

	room.Items.Update(item)
}
