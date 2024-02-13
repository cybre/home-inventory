package eventsourcing

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cybre/home-inventory/pkg/logging"
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

type EventPublisher interface {
	PublishEvents(ctx context.Context, events []Event) error
}

type CommandBus struct {
	eventStore     EventStore
	eventPublisher EventPublisher
}

func NewCommandBus(eventStore EventStore, eventPublisher EventPublisher) *CommandBus {
	return &CommandBus{
		eventStore:     eventStore,
		eventPublisher: eventPublisher,
	}
}

func (cb *CommandBus) Dispatch(ctx context.Context, c Command) error {
	// TODO - distributed locking
	events, err := cb.eventStore.GetEvents(c.AggregateType(), c.AggregateID())
	if err != nil {
		return fmt.Errorf("failed to fetch events for aggregate: %w", err)
	}

	aggregateVersion := uint(0)
	if len(events) > 0 {
		aggregateVersion = events[len(events)-1].Version
	}

	aggregateContext := NewAggregateContext(c.AggregateType(), c.AggregateID(), aggregateVersion)
	aggregate, ok := GetAggregateRoot(aggregateContext)
	if !ok {
		return ErrAggregateTypeNotFound
	}

	for _, event := range events {
		aggregate.ApplyEvent(event.Data)
	}

	aggCtx := logging.WithLogger(
		ctx,
		logging.FromContext(ctx).With(
			slog.Any("aggregate_type", c.AggregateType()),
			slog.Any("aggregate_id", c.AggregateID()),
			slog.Any("version", aggregateContext.Version()),
			slog.String("command", fmt.Sprintf("%T", c)),
		),
	)

	result, err := aggregate.HandleCommand(aggCtx, c)
	if err != nil {
		return fmt.Errorf("aggregate failed to handle command: %w", err)
	}

	newEvents := utils.Map(result, func(i uint, event EventData) Event {
		return Event{
			AggregateType: c.AggregateType(),
			AggregateID:   c.AggregateID(),
			EventType:     event.EventType(),
			Data:          event,
			Timestamp:     time.Now().UnixMilli(),
			Version:       aggregateContext.Version() + i + 1,
		}
	})

	if err := cb.eventStore.StoreEvents(ctx, newEvents); err != nil {
		return fmt.Errorf("failed to store events: %w", err)
	}

	if err := cb.eventPublisher.PublishEvents(ctx, newEvents); err != nil {
		return fmt.Errorf("failed to publish events: %w", err)
	}

	return nil
}
