package user

import es "github.com/cybre/home-inventory/pkg/eventsourcing"

const (
	EventTypeUserCreated             es.EventType = "UserCreatedEvent"
	EventTypeUserLoginTokenGenerated es.EventType = "UserLoginTokenGeneratedEvent"
	EventTypeUserLogin               es.EventType = "UserLoginEvent"
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

type UserLoginTokenGeneratedEvent struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

func (e UserLoginTokenGeneratedEvent) EventType() es.EventType {
	return EventTypeUserLoginTokenGenerated
}

type UserLoginEvent struct {
	UserID    string `json:"userId"`
	Timestamp int64  `json:"timestamp"`
	UserAgent string `json:"userAgent"`
	IP        string `json:"ip"`
}

func (e UserLoginEvent) EventType() es.EventType {
	return EventTypeUserLogin
}
