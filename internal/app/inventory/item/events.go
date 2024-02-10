package item

import "github.com/cybre/home-inventory/pkg/domain"

const (
	EventTypeItemAdded   domain.EventType = "ItemAddedEvent"
	EventTypeItemUpdated domain.EventType = "ItemUpdatedEvent"
)

type ItemAddedEvent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Barcode  string `json:"barcode"`
	Quantity uint   `json:"quantity"`
}

func (e ItemAddedEvent) EventType() domain.EventType {
	return EventTypeItemAdded
}

type ItemUpdatedEvent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Barcode  string `json:"barcode"`
	Quantity uint   `json:"quantity"`
}

func (e ItemUpdatedEvent) EventType() domain.EventType {
	return EventTypeItemUpdated
}
