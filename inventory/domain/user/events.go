package user

import es "github.com/cybre/home-inventory/pkg/eventsourcing"

const (
	EventTypeUserCreated               es.EventType = "UserCreatedEvent"
	EventTypeUserOneTimeTokenGenerated es.EventType = "UserOneTimeTokenGeneratedEvent"
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

type UserOneTimeTokenGeneratedEvent struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

func (e UserOneTimeTokenGeneratedEvent) EventType() es.EventType {
	return EventTypeUserOneTimeTokenGenerated
}
