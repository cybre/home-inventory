package item

import (
	"context"
	"fmt"

	"github.com/cybre/home-inventory/pkg/domain"
)

type ItemProjector struct {
}

func NewItemProjector() *ItemProjector {
	return &ItemProjector{}
}

func (p ItemProjector) HandleEvent(ctx context.Context, event domain.EventData) error {
	switch e := event.(type) {
	case ItemAddedEvent:
		return p.handleItemAddedEvent(ctx, e)
	case ItemUpdatedEvent:
		return p.handleItemUpdatedEvent(ctx, e)
	default:
		return domain.ErrUnknownEvent
	}
}

func (p ItemProjector) Events() []domain.EventType {
	return []domain.EventType{
		EventTypeItemAdded,
		EventTypeItemUpdated,
	}
}

func (p ItemProjector) Name() string {
	return "item.Projector"
}

func (p ItemProjector) handleItemAddedEvent(ctx context.Context, e ItemAddedEvent) error {
	fmt.Printf("item added: %#v\n", e)
	return nil
}

func (p ItemProjector) handleItemUpdatedEvent(ctx context.Context, e ItemUpdatedEvent) error {
	fmt.Printf("item updated: %#v\n", e)
	return nil
}