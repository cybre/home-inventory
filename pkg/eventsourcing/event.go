package eventsourcing

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type EventType string

type Event struct {
	AggregateType AggregateType `json:"aggregateType"`
	AggregateID   AggregateID   `json:"aggregateId"`
	EventType     EventType     `json:"eventType"`
	Data          EventData     `json:"eventData"`
	Timestamp     int64         `json:"timestamp"`
	Version       uint          `json:"version"`
}

func UnmarshalEvent(data []byte) (Event, error) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		return Event{}, fmt.Errorf("failed to decode event: %w", err)
	}

	eventDataInstance, ok := GetEvent(EventType(event["eventType"].(string)))
	if !ok {
		return Event{}, ErrEventTypeNotFound
	}

	eventDataBytes, err := json.Marshal(event["eventData"])
	if err != nil {
		return Event{}, fmt.Errorf("failed to marshal event data: %w", err)
	}

	if err := json.Unmarshal(eventDataBytes, eventDataInstance); err != nil {
		return Event{}, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	return Event{
		AggregateType: AggregateType(event["aggregateType"].(string)),
		AggregateID:   AggregateID(event["aggregateId"].(string)),
		EventType:     EventType(event["eventType"].(string)),
		Data:          reflect.ValueOf(eventDataInstance).Elem().Interface().(EventData),
		Timestamp:     int64(event["timestamp"].(float64)),
		Version:       uint(event["version"].(float64)),
	}, nil
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
