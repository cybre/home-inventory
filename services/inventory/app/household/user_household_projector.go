package household

import (
	"context"
	"fmt"

	es "github.com/cybre/home-inventory/internal/eventsourcing"
	"github.com/cybre/home-inventory/services/inventory/domain/household"
)

type HouseholdRepo interface {
	InsertHousehold(ctx context.Context, userId string, householdId string, name string, location string, description string) error
	UpdateHousehold(ctx context.Context, userId string, householdId string, name string, location string, description string) error

	UpsertRoom(ctx context.Context, userId string, householdId string, roomId string, name string, itemCount int) error
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
	if err := p.repository.InsertHousehold(ctx, e.UserID, e.HouseholdID, e.Name, e.Location, e.Description); err != nil {
		return fmt.Errorf("failed to insert household: %w", err)
	}

	return nil
}

func (p UserHouseholdProjector) handleHouseholdUpdatedEvent(ctx context.Context, e household.HouseholdUpdatedEvent) error {
	if err := p.repository.UpdateHousehold(ctx, e.UserID, e.HouseholdID, e.Name, e.Location, e.Description); err != nil {
		return fmt.Errorf("failed to update household: %w", err)
	}

	return nil
}

func (p UserHouseholdProjector) handleRoomAddedEvent(ctx context.Context, e household.RoomAddedEvent) error {
	if err := p.repository.UpsertRoom(ctx, e.UserID, e.HouseholdID, e.RoomID, e.Name, 0); err != nil {
		return fmt.Errorf("failed to add room: %w", err)
	}

	return nil
}

func (p UserHouseholdProjector) handleRoomUpdatedEvent(ctx context.Context, e household.RoomUpdatedEvent) error {
	if err := p.repository.UpsertRoom(ctx, e.UserID, e.HouseholdID, e.RoomID, e.Name, e.ItemCount); err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	return nil
}
