package household

import (
	"context"
	"log/slog"

	es "github.com/cybre/home-inventory/internal/eventsourcing"
	"github.com/cybre/home-inventory/internal/logging"
	"github.com/cybre/home-inventory/services/inventory/domain/household"
)

type ItemProjector struct {
}

func NewItemProjector() *ItemProjector {
	return &ItemProjector{}
}

func (p ItemProjector) HandleEvent(ctx context.Context, event es.EventData) error {
	switch e := event.(type) {
	case household.ItemAddedEvent:
		return p.handleItemAddedEvent(ctx, e)
	case household.ItemUpdatedEvent:
		return p.handleItemUpdatedEvent(ctx, e)
	default:
		return es.ErrUnknownEvent
	}
}

func (p ItemProjector) Events() []es.EventType {
	return []es.EventType{
		household.EventTypeItemAdded,
		household.EventTypeItemUpdated,
	}
}

func (p ItemProjector) Name() string {
	return "household.ItemProjector"
}

func (p ItemProjector) handleItemAddedEvent(ctx context.Context, e household.ItemAddedEvent) error {
	logging.FromContext(ctx).Info("item added", slog.String("name", e.Name))

	return nil
}

func (p ItemProjector) handleItemUpdatedEvent(ctx context.Context, e household.ItemUpdatedEvent) error {
	logging.FromContext(ctx).Info("item updated", slog.String("name", e.Name))

	return nil
}
