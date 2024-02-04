package domain

import (
	"context"
	"time"

	"github.com/cybre/home-inventory/pkg/utils"
)

type Command interface {
	AggregateType() AggregateType
	AggregateID() AggregateID
}

type EventStore interface {
	GetEvents(aggregateType AggregateType, aggregateID AggregateID) ([]Event, error)
	StoreEvents(ctx context.Context, events []Event) error
}

type CommandBus struct {
	eventStore EventStore
}

func NewCommandBus(eventStore EventStore) *CommandBus {
	return &CommandBus{eventStore: eventStore}
}

func (h *CommandBus) Dispatch(ctx context.Context, c Command) error {
	// TODO - distributed locking
	events, err := h.eventStore.GetEvents(c.AggregateType(), c.AggregateID())
	if err != nil {
		return err
	}

	aggregateVersion := uint(0)
	if len(events) > 0 {
		aggregateVersion = events[len(events)-1].version
	}

	aggregateContext := NewAggregateContext(c.AggregateType(), c.AggregateID(), aggregateVersion)
	aggregate, ok := GetAggregateRoot(aggregateContext)
	if !ok {
		return ErrAggregateTypeNotFound
	}

	for _, event := range events {
		aggregate.ApplyEvent(event.eventData)
	}

	result, err := aggregate.HandleCommand(ctx, c)
	if err != nil {
		return err
	}

	return h.eventStore.StoreEvents(ctx, utils.Map(result, func(i uint, event EventData) Event {
		return Event{
			aggregateType: c.AggregateType(),
			aggregateID:   c.AggregateID(),
			eventType:     event.EventType(),
			eventData:     event,
			timestamp:     time.Now().UnixMilli(),
			version:       aggregateContext.Version() + i + 1,
		}
	}))
}
