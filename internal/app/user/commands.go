package user

import (
	"github.com/cybre/home-inventory/pkg/domain"
)

type CreateUserCommand struct {
	UserID    string
	FirstName string
	LastName  string
}

func (c CreateUserCommand) AggregateType() domain.AggregateType {
	return UserAggregateType
}

func (c CreateUserCommand) AggregateID() domain.AggregateID {
	return domain.AggregateID(c.UserID)
}
