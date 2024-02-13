package user

import es "github.com/cybre/home-inventory/pkg/eventsourcing"

const (
	EventTypeUserCreated es.EventType = "UserCreatedEvent"
)

type UserCreatedEvent struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (e UserCreatedEvent) EventType() es.EventType {
	return EventTypeUserCreated
}
