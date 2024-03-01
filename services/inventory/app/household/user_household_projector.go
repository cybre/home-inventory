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
	default:
		return es.ErrUnknownEvent
	}
}

func (p UserHouseholdProjector) Events() []es.EventType {
	return []es.EventType{
		household.EventTypeHouseholdCreated,
		household.EventTypeHouseholdUpdated,
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
		return fmt.Errorf("failed to insert household: %w", err)
	}

	return nil
}
