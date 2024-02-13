package user

import (
	"context"
	"fmt"

	c "github.com/cybre/home-inventory/inventory/domain/common"
	es "github.com/cybre/home-inventory/pkg/eventsourcing"
)

const (
	UserAggregateType       es.AggregateType = "UserAggregate"
	initialAggregateVersion                  = 0
)

type UserAggregate struct {
	es.AggregateContext

	ID        c.UserID
	FirstName FirstName
	LastName  LastName
	Email     Email
}

func NewHouseholdAggregate(aggregateContext es.AggregateContext) es.AggregateRoot {
	return &UserAggregate{
		AggregateContext: aggregateContext,
	}
}

func (a *UserAggregate) ApplyEvent(event es.EventData) {
	switch e := event.(type) {
	case UserCreatedEvent:
		a.applyUserCreatedEvent(e)
	default:
		panic("unknown event type")
	}
}

func (a *UserAggregate) HandleCommand(ctx context.Context, command es.Command) ([]es.EventData, error) {
	switch c := command.(type) {
	case CreateUserCommand:
		return a.handleCreateUserCommand(ctx, c)
	default:
		return nil, es.ErrUnknownCommand
	}
}

func (a *UserAggregate) handleCreateUserCommand(ctx context.Context, command CreateUserCommand) ([]es.EventData, error) {
	if a.Version() != initialAggregateVersion {
		return nil, fmt.Errorf("user with provided ID already exists: %s", command.UserID)
	}

	userId, err := c.NewUserID(command.UserID)
	if err != nil {
		return nil, err
	}

	firstName, err := NewFirstName(command.FirstName)
	if err != nil {
		return nil, err
	}

	lastName, err := NewLastName(command.LastName)
	if err != nil {
		return nil, err
	}

	email, err := NewEmail(command.Email)
	if err != nil {
		return nil, err
	}

	return c.Events(UserCreatedEvent{
		UserID:    userId.String(),
		FirstName: firstName.String(),
		LastName:  lastName.String(),
		Email:     email.String(),
	})
}

func (a *UserAggregate) applyUserCreatedEvent(event UserCreatedEvent) {
	a.ID, _ = c.NewUserID(event.UserID)
	a.FirstName, _ = NewFirstName(event.FirstName)
	a.LastName, _ = NewLastName(event.LastName)
}
