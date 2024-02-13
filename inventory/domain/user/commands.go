package user

import es "github.com/cybre/home-inventory/pkg/eventsourcing"

type CreateUserCommand struct {
	UserID    string
	FirstName string
	LastName  string
	Email     string
}

func (c CreateUserCommand) AggregateType() es.AggregateType {
	return UserAggregateType
}

func (c CreateUserCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.UserID)
}
