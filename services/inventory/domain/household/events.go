package household

import es "github.com/cybre/home-inventory/internal/eventsourcing"

const (
	EventTypeHouseholdCreated es.EventType = "HouseholdCreatedEvent"
	EventTypeHouseholdUpdated es.EventType = "HouseholdUpdatedEvent"
	EventTypeRoomAdded        es.EventType = "RoomAddedEvent"
	EventTypeItemAdded        es.EventType = "ItemAddedEvent"
	EventTypeItemUpdated      es.EventType = "ItemUpdatedEvent"
)

type HouseholdCreatedEvent struct {
	HouseholdID string `json:"householdId"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
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
}

func (e HouseholdUpdatedEvent) EventType() es.EventType {
	return EventTypeHouseholdUpdated
}

type RoomAddedEvent struct {
	HouseholdID string `json:"householdId"`
	RoomID      string `json:"roomId"`
	Name        string `json:"name"`
}

func (e RoomAddedEvent) EventType() es.EventType {
	return EventTypeRoomAdded
}

type ItemAddedEvent struct {
	HouseholdID string `json:"householdId"`
	RoomID      string `json:"roomId"`
	ItemID      string `json:"itemId"`
	Name        string `json:"name"`
	Barcode     string `json:"barcode"`
	Quantity    uint   `json:"quantity"`
}

func (e ItemAddedEvent) EventType() es.EventType {
	return EventTypeItemAdded
}

type ItemUpdatedEvent struct {
	HouseholdID string `json:"householdId"`
	RoomID      string `json:"roomId"`
	ItemID      string `json:"itemId"`
	Name        string `json:"name"`
	Barcode     string `json:"barcode"`
	Quantity    uint   `json:"quantity"`
}

func (e ItemUpdatedEvent) EventType() es.EventType {
	return EventTypeItemUpdated
}
