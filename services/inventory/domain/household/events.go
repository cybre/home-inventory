package household

import es "github.com/cybre/home-inventory/internal/eventsourcing"

const (
	EventTypeHouseholdCreated es.EventType = "HouseholdCreatedEvent"
	EventTypeHouseholdUpdated es.EventType = "HouseholdUpdatedEvent"
	EventTypeHouseholdDeleted es.EventType = "HouseholdDeletedEvent"

	EventTypeRoomAdded   es.EventType = "RoomAddedEvent"
	EventTypeRoomUpdated es.EventType = "RoomUpdatedEvent"
	EventTypeRoomDeleted es.EventType = "RoomDeletedEvent"
)

type HouseholdCreatedEvent struct {
	HouseholdID string `json:"householdId"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Order       uint   `json:"order"`
	Timestamp   int64  `json:"timestamp"`
}

func (e HouseholdCreatedEvent) EventType() es.EventType {
	return EventTypeHouseholdCreated
}

type HouseholdUpdatedEvent struct {
	HouseholdID string `json:"householdId"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Timestamp   int64  `json:"timestamp"`
}

func (e HouseholdUpdatedEvent) EventType() es.EventType {
	return EventTypeHouseholdUpdated
}

type HouseholdDeletedEvent struct {
	HouseholdID string `json:"householdId"`
	UserID      string `json:"userId"`
}

func (e HouseholdDeletedEvent) EventType() es.EventType {
	return EventTypeHouseholdDeleted
}

type RoomAddedEvent struct {
	HouseholdID string `json:"householdId"`
	UserID      string `json:"userId"`
	RoomID      string `json:"roomId"`
	Name        string `json:"name"`
	Order       uint   `json:"order"`
	Timestamp   int64  `json:"timestamp"`
}

func (e RoomAddedEvent) EventType() es.EventType {
	return EventTypeRoomAdded
}

type RoomUpdatedEvent struct {
	HouseholdID string `json:"householdId"`
	UserID      string `json:"userId"`
	RoomID      string `json:"roomId"`
	Name        string `json:"name"`
	Order       uint   `json:"order"`
	Timestamp   int64  `json:"timestamp"`
}

func (e RoomUpdatedEvent) EventType() es.EventType {
	return EventTypeRoomUpdated
}

type RoomDeletedEvent struct {
	HouseholdID string `json:"householdId"`
	UserID      string `json:"userId"`
	RoomID      string `json:"roomId"`
}

func (e RoomDeletedEvent) EventType() es.EventType {
	return EventTypeRoomDeleted
}
