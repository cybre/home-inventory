package user

import "github.com/cybre/home-inventory/pkg/domain"

const (
	EventTypeUserCreated domain.EventType = "UserCreatedEvent"
)

type UserCreatedEvent struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (e UserCreatedEvent) EventType() domain.EventType {
	return EventTypeUserCreated
}
