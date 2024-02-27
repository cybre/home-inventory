package household

import (
	"context"
	"fmt"

	es "github.com/cybre/home-inventory/internal/eventsourcing"
	"github.com/cybre/home-inventory/services/inventory/domain/household"
)

type HouseholdInserter interface {
	InsertHousehold(ctx context.Context, userId string, householdId string, name string, location string, description string) error
}

type UserHouseholdProjector struct {
	repository HouseholdInserter
}

func NewUserHouseholdProjector(repository HouseholdInserter) *UserHouseholdProjector {
	return &UserHouseholdProjector{
		repository: repository,
	}
}

func (p UserHouseholdProjector) HandleEvent(ctx context.Context, event es.EventData) error {
	switch e := event.(type) {
	case household.HouseholdCreatedEvent:
		return p.handleHouseholdCreatedEvent(ctx, e)
	default:
		return es.ErrUnknownEvent
	}
}

func (p UserHouseholdProjector) Events() []es.EventType {
	return []es.EventType{
		household.EventTypeHouseholdCreated,
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
