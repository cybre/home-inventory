package user

import "github.com/cybre/home-inventory/pkg/domain"

const (
	EventTypeUserCreated domain.EventType = "UserCreatedEvent"
)

type UserCreatedEvent struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (e UserCreatedEvent) EventType() domain.EventType {
	return EventTypeUserCreated
}
