package household

import "github.com/cybre/home-inventory/pkg/domain"

const (
	EventTypeHouseholdCreated domain.EventType = "HouseholdCreatedEvent"
	EventTypeRoomAdded        domain.EventType = "RoomAddedEvent"
	EventTypeItemAdded        domain.EventType = "ItemAddedEvent"
	EventTypeItemUpdated      domain.EventType = "ItemUpdatedEvent"
)

type HouseholdCreatedEvent struct {
	HouseholdID string `json:"householdId"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
}

func (e HouseholdCreatedEvent) EventType() domain.EventType {
	return EventTypeHouseholdCreated
}

type RoomAddedEvent struct {
	HouseholdID string `json:"householdId"`
	RoomID      string `json:"roomId"`
	Name        string `json:"name"`
}

func (e RoomAddedEvent) EventType() domain.EventType {
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

func (e ItemAddedEvent) EventType() domain.EventType {
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

func (e ItemUpdatedEvent) EventType() domain.EventType {
	return EventTypeItemUpdated
}
