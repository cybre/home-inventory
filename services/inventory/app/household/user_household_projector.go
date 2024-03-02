package household

import (
	"context"
	"fmt"

	es "github.com/cybre/home-inventory/internal/eventsourcing"
	"github.com/cybre/home-inventory/services/inventory/domain/household"
	"github.com/gocql/gocql"
)

type HouseholdRepo interface {
	InsertHousehold(ctx context.Context, model UserHouseholdModel) error
	UpdateHousehold(ctx context.Context, model UserHouseholdModel) error

	UpsertRoom(ctx context.Context, userId string, model UserHouseholdRoomModel) error
}

type UserHouseholdProjector struct {
	repository HouseholdRepo
}

func NewUserHouseholdProjector(repository HouseholdRepo) *UserHouseholdProjector {
	return &UserHouseholdProjector{
		repository: repository,
	}
}

func (p UserHouseholdProjector) HandleEvent(ctx context.Context, event es.EventData) error {
	switch e := event.(type) {
	case household.HouseholdCreatedEvent:
		return p.handleHouseholdCreatedEvent(ctx, e)
	case household.HouseholdUpdatedEvent:
		return p.handleHouseholdUpdatedEvent(ctx, e)
	case household.RoomAddedEvent:
		return p.handleRoomAddedEvent(ctx, e)
	case household.RoomUpdatedEvent:
		return p.handleRoomUpdatedEvent(ctx, e)
	default:
		return es.ErrUnknownEvent
	}
}

func (p UserHouseholdProjector) Events() []es.EventType {
	return []es.EventType{
		household.EventTypeHouseholdCreated,
		household.EventTypeHouseholdUpdated,
		household.EventTypeRoomAdded,
		household.EventTypeRoomUpdated,
	}
}

func (p UserHouseholdProjector) Name() string {
	return "household.UserHouseholdProjector"
}

func (p UserHouseholdProjector) handleHouseholdCreatedEvent(ctx context.Context, e household.HouseholdCreatedEvent) error {
	householdUUID, err := gocql.ParseUUID(e.HouseholdID)
	if err != nil {
		return fmt.Errorf("failed to parse household ID: %w", err)
	}

	if err := p.repository.InsertHousehold(ctx, UserHouseholdModel{
		UserID:      e.UserID,
		HouseholdID: householdUUID,
		Name:        e.Name,
		Location:    e.Location,
		Description: e.Description,
		Order:       e.Order,
		Timestamp:   e.Timestamp,
	}); err != nil {
		return fmt.Errorf("failed to insert household: %w", err)
	}

	return nil
}

func (p UserHouseholdProjector) handleHouseholdUpdatedEvent(ctx context.Context, e household.HouseholdUpdatedEvent) error {
	householdUUID, err := gocql.ParseUUID(e.HouseholdID)
	if err != nil {
		return fmt.Errorf("failed to parse household ID: %w", err)
	}

	if err := p.repository.UpdateHousehold(ctx, UserHouseholdModel{
		UserID:      e.UserID,
		HouseholdID: householdUUID,
		Name:        e.Name,
		Location:    e.Location,
		Description: e.Description,
		Timestamp:   e.Timestamp,
	}); err != nil {
		return fmt.Errorf("failed to update household: %w", err)
	}

	return nil
}

func (p UserHouseholdProjector) handleRoomAddedEvent(ctx context.Context, e household.RoomAddedEvent) error {
	householdUUID, err := gocql.ParseUUID(e.HouseholdID)
	if err != nil {
		return fmt.Errorf("failed to parse household ID: %w", err)
	}

	roomUUID, err := gocql.ParseUUID(e.RoomID)
	if err != nil {
		return fmt.Errorf("failed to parse room ID: %w", err)
	}

	if err := p.repository.UpsertRoom(ctx, e.UserID, UserHouseholdRoomModel{
		HouseholdID: householdUUID,
		RoomID:      roomUUID,
		Name:        e.Name,
		Order:       e.Order,
		Timestamp:   e.Timestamp,
	}); err != nil {
		return fmt.Errorf("failed to add room: %w", err)
	}

	return nil
}

func (p UserHouseholdProjector) handleRoomUpdatedEvent(ctx context.Context, e household.RoomUpdatedEvent) error {
	householdUUID, err := gocql.ParseUUID(e.HouseholdID)
	if err != nil {
		return fmt.Errorf("failed to parse household ID: %w", err)
	}

	roomUUID, err := gocql.ParseUUID(e.RoomID)
	if err != nil {
		return fmt.Errorf("failed to parse room ID: %w", err)
	}

	if err := p.repository.UpsertRoom(ctx, e.UserID, UserHouseholdRoomModel{
		HouseholdID: householdUUID,
		RoomID:      roomUUID,
		Name:        e.Name,
		Order:       e.Order,
		Timestamp:   e.Timestamp,
	}); err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	return nil
}
