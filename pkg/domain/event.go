package domain

import (
	"encoding/json"
	"reflect"
)

type EventType string

type Event struct {
	AggregateType AggregateType `json:"aggregate_type"`
	AggregateID   AggregateID   `json:"aggregate_id"`
	EventType     EventType     `json:"event_type"`
	Data          EventData     `json:"event_data"`
	Timestamp     int64         `json:"timestamp"`
	Version       uint          `json:"version"`
}

func (e Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
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
