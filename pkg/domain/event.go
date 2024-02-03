package domain

import "reflect"

// EventType represents the type of an event.
type EventType string

// Event represents an event.
type Event struct {
	aggregateType AggregateType
	aggregateID   AggregateID
	eventType     EventType
	eventData     EventData
	timestamp     int64
	version       uint
}

func NewEvent(aggregateType AggregateType, aggregateID AggregateID, eventType EventType, eventData EventData, timestamp int64, version uint) Event {
	return Event{
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
		eventType:     eventType,
		eventData:     eventData,
		timestamp:     timestamp,
		version:       version,
	}
}

func (e Event) AggregateType() AggregateType {
	return e.aggregateType
}

func (e Event) AggregateID() AggregateID {
	return e.aggregateID
}

func (e Event) EventType() EventType {
	return e.eventType
}

func (e Event) EventData() EventData {
	return e.eventData
}

func (e Event) Timestamp() int64 {
	return e.timestamp
}

func (e Event) Version() uint {
	return e.version
}

type EventData interface {
	EventType() EventType
}

type EventRegistry struct {
	events map[EventType]reflect.Type
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		events: make(map[EventType]reflect.Type),
	}
}

func (r *EventRegistry) RegisterEvent(event EventData) {
	r.events[event.EventType()] = reflect.TypeOf(event)
}

func (r *EventRegistry) GetEvent(eventType EventType) (EventData, bool) {
	eventReflectType, ok := r.events[eventType]
	if ok {
		return reflect.New(eventReflectType).Interface().(EventData), true
	}

	return nil, false
}

var eventRegistry *EventRegistry

func init() {
	eventRegistry = NewEventRegistry()
}

func RegisterEvent(event EventData) {
	eventRegistry.RegisterEvent(event)
}

func GetEvent(eventType EventType) (EventData, bool) {
	return eventRegistry.GetEvent(eventType)
}
