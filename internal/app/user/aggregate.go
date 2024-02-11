package user

import (
	"context"
	"fmt"

	c "github.com/cybre/home-inventory/internal/app/common"
	"github.com/cybre/home-inventory/pkg/domain"
)

const (
	UserAggregateType       domain.AggregateType = "UserAggregate"
	initialAggregateVersion                      = 0
)

type UserAggregate struct {
	domain.AggregateContext

	ID        c.UserID
	FirstName FirstName
	LastName  LastName
}

func NewHouseholdAggregate(aggregateContext domain.AggregateContext) domain.AggregateRoot {
	return &UserAggregate{
		AggregateContext: aggregateContext,
	}
}

func (a *UserAggregate) ApplyEvent(event domain.EventData) {
	switch e := event.(type) {
	case UserCreatedEvent:
		a.applyUserCreatedEvent(e)
	default:
		panic("unknown event type")
	}
}

func (a *UserAggregate) HandleCommand(ctx context.Context, command domain.Command) ([]domain.EventData, error) {
	switch c := command.(type) {
	case CreateUserCommand:
		return a.handleCreateUserCommand(ctx, c)
	default:
		return nil, domain.ErrUnknownCommand
	}
}

func (a *UserAggregate) handleCreateUserCommand(ctx context.Context, command CreateUserCommand) ([]domain.EventData, error) {
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

	return c.Events(UserCreatedEvent{
		ID:        userId.String(),
		FirstName: firstName.String(),
		LastName:  lastName.String(),
	})
}

func (a *UserAggregate) applyUserCreatedEvent(event UserCreatedEvent) {
	a.ID, _ = c.NewUserID(event.ID)
	a.FirstName, _ = NewFirstName(event.FirstName)
	a.LastName, _ = NewLastName(event.LastName)
}
