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

type GenerateLoginTokenCommand struct {
	UserID string
}

func (c GenerateLoginTokenCommand) AggregateType() es.AggregateType {
	return UserAggregateType
}

func (c GenerateLoginTokenCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.UserID)
}

type LoginCommand struct {
	UserID    string
	UserAgent string
	IP        string
}

func (c LoginCommand) AggregateType() es.AggregateType {
	return UserAggregateType
}

func (c LoginCommand) AggregateID() es.AggregateID {
	return es.AggregateID(c.UserID)
}
