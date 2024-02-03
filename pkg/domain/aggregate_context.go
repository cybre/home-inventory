package domain

import (
	"time"
)

type AggregateContext struct {
	aggregateType AggregateType
	aggregateID   AggregateID
	version       uint

	events []Event
}

func NewAggregateContext(aggregateType AggregateType, aggregateID AggregateID, version uint) AggregateContext {
	return AggregateContext{
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
		version:       version,
		events:        []Event{},
	}
}

func (a *AggregateContext) AggregateID() AggregateID {
	return a.aggregateID
}

func (a *AggregateContext) AggregateType() AggregateType {
	return a.aggregateType
}

func (a *AggregateContext) StoreEvent(event EventData) {
	a.version++

	a.events = append(a.events, Event{
		aggregateType: a.aggregateType,
		aggregateID:   a.aggregateID,
		eventType:     event.EventType(),
		eventData:     event,
		timestamp:     time.Now().UnixMilli(),
		version:       a.version,
	})
}

func (a *AggregateContext) Events() []Event {
	return a.events
}

func (a *AggregateContext) Version() uint {
	return a.version
}
