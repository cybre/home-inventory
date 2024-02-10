package items

import "github.com/cybre/home-inventory/pkg/domain"

const (
	EventTypeItemAdded   domain.EventType = "ItemAddedEvent"
	EventTypeItemUpdated domain.EventType = "ItemUpdatedEvent"
)

type ItemAddedEvent struct {
	Name string `json:"name"`
}

func (e ItemAddedEvent) EventType() domain.EventType {
	return EventTypeItemAdded
}

type ItemUpdatedEvent struct {
	Name string `json:"name"`
}

func (e ItemUpdatedEvent) EventType() domain.EventType {
	return EventTypeItemUpdated
}
