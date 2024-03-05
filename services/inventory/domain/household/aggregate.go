package household

import (
	"context"
	"time"

	"github.com/bnkamalesh/errors"
	es "github.com/cybre/home-inventory/internal/eventsourcing"
	c "github.com/cybre/home-inventory/services/inventory/domain/common"
)

const (
	HouseholdAggregateType  es.AggregateType = "HouseholdAggregate"
	initialAggregateVersion                  = 0
)

type HouseholdAgregate struct {
	es.AggregateContext

	UserID      c.UserID
	Name        HouseholdName
	Location    HouseholdLocation
	Description HouseholdDescription
	Order       uint

	Rooms Rooms

	Deleted bool
}

func NewHouseholdAggregate(aggregateContext es.AggregateContext) es.AggregateRoot {
	return &HouseholdAgregate{
		AggregateContext: aggregateContext,
	}
}

func (a *HouseholdAgregate) ApplyEvent(event es.EventData) {
	switch e := event.(type) {
	case HouseholdCreatedEvent:
		a.applyHouseholdCreatedEvent(e)
	case HouseholdUpdatedEvent:
		a.applyHouseholdUpdatedEvent(e)
	case HouseholdDeletedEvent:
		a.applyHouseholdDeletedEvent(e)
	case RoomAddedEvent:
		a.applyRoomAddedEvent(e)
	case RoomUpdatedEvent:
		a.applyRoomUpdatedEvent(e)
	case RoomDeletedEvent:
		a.applyRoomDeletedEvent(e)
	default:
		panic("unknown event type")
	}
}

func (a *HouseholdAgregate) HandleCommand(ctx context.Context, command es.Command) ([]es.EventData, error) {
	if _, ok := command.(CreateHouseholdCommand); !ok {
		if a.Version() == initialAggregateVersion || a.Deleted {
			return nil, errors.NotFound("household with provided ID does not exist")
		}
	} else if a.Version() != initialAggregateVersion {
		return nil, errors.Duplicate("household with provided ID already exists")
	}

	switch c := command.(type) {
	case CreateHouseholdCommand:
		return a.handleCreateHouseholdCommand(ctx, c)
	case UpdateHouseholdCommand:
		return a.handleUpdateHouseholdCommand(ctx, c)
	case DeleteHouseholdCommand:
		return a.handleDeleteHouseholdCommand(ctx, c)
	case AddRoomCommand:
		return a.handleAddRoomCommand(ctx, c)
	case UpdateRoomCommand:
		return a.handleUpdateRoomCommand(ctx, c)
	case DeleteRoomCommand:
		return a.handleDeleteRoomCommand(ctx, c)
	default:
		return nil, es.ErrUnknownCommand
	}
}

func (a *HouseholdAgregate) handleCreateHouseholdCommand(ctx context.Context, command CreateHouseholdCommand) ([]es.EventData, error) {
	userId, err := c.NewUserID(command.UserID)
	if err != nil {
		return nil, err
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

	return c.Events(HouseholdCreatedEvent{
		HouseholdID: a.AggregateID().String(),
		UserID:      userId.String(),
		Name:        name.String(),
		Location:    location.String(),
		Description: description.String(),
		Order:       command.Order,
		Timestamp:   time.Now().UnixMilli(),
	})
}

func (a *HouseholdAgregate) handleUpdateHouseholdCommand(ctx context.Context, command UpdateHouseholdCommand) ([]es.EventData, error) {
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

func (a *HouseholdAgregate) handleDeleteHouseholdCommand(ctx context.Context, command DeleteHouseholdCommand) ([]es.EventData, error) {
	return c.Events(HouseholdDeletedEvent{
		HouseholdID: a.AggregateID().String(),
		UserID:      a.UserID.String(),
	})
}

func (a *HouseholdAgregate) handleAddRoomCommand(ctx context.Context, command AddRoomCommand) ([]es.EventData, error) {
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
		Order:       newRoom.Order,
		Timestamp:   time.Now().UnixMilli(),
	})
}

func (a *HouseholdAgregate) handleUpdateRoomCommand(ctx context.Context, command UpdateRoomCommand) ([]es.EventData, error) {
	roomID, err := NewRoomID(command.RoomID)
	if err != nil {
		return nil, err
	}

	room, ok := a.Rooms.Get(roomID)
	if !ok {
		return nil, errors.NotFoundf("room with ID %s does not exist", roomID)
	}

	room, err = room.Update(command.Name)
	if err != nil {
		return nil, err
	}

	if err := a.Rooms.Update(room); err != nil {
		return nil, err
	}

	return c.Events(RoomUpdatedEvent{
		HouseholdID: a.AggregateID().String(),
		UserID:      a.UserID.String(),
		RoomID:      room.ID.String(),
		Name:        room.Name.String(),
		Order:       room.Order,
		Timestamp:   time.Now().UnixMilli(),
	})
}

func (a *HouseholdAgregate) handleDeleteRoomCommand(ctx context.Context, command DeleteRoomCommand) ([]es.EventData, error) {
	roomID, err := NewRoomID(command.RoomID)
	if err != nil {
		return nil, err
	}

	return c.Events(RoomDeletedEvent{
		HouseholdID: a.AggregateID().String(),
		UserID:      a.UserID.String(),
		RoomID:      roomID.String(),
	})
}

func (a *HouseholdAgregate) applyHouseholdCreatedEvent(event HouseholdCreatedEvent) {
	a.UserID, _ = c.NewUserID(event.UserID)
	a.Name, _ = NewHouseholdName(event.Name)
	a.Location, _ = NewHouseholdLocation(event.Location)
	a.Description, _ = NewHouseholdDescription(event.Description)
	a.Order = event.Order
	a.Rooms = NewRooms()
}

func (a *HouseholdAgregate) applyHouseholdUpdatedEvent(event HouseholdUpdatedEvent) {
	a.Name, _ = NewHouseholdName(event.Name)
	a.Location, _ = NewHouseholdLocation(event.Location)
	a.Description, _ = NewHouseholdDescription(event.Description)
}

func (a *HouseholdAgregate) applyHouseholdDeletedEvent(event HouseholdDeletedEvent) {
	a.Deleted = true
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

func (a *HouseholdAgregate) applyRoomDeletedEvent(event RoomDeletedEvent) {
	roomID, _ := NewRoomID(event.RoomID)
	a.Rooms.Remove(roomID)
}
