package user

import (
	"context"
	"fmt"
	"time"

	c "github.com/cybre/home-inventory/inventory/domain/common"
	es "github.com/cybre/home-inventory/pkg/eventsourcing"
	"github.com/google/uuid"
)

type EmailUniquenessService interface {
	IsUnique(ctx context.Context, email Email) (bool, error)
}

const (
	UserAggregateType       es.AggregateType = "UserAggregate"
	initialAggregateVersion                  = 0
)

type UserAggregate struct {
	es.AggregateContext

	emailUniquenessService EmailUniquenessService

	ID        c.UserID
	FirstName FirstName
	LastName  LastName
	Email     Email

	LastLogin LastLogin
}

func NewUserAggregate(emailUniquenessService EmailUniquenessService) es.AggregateRootFactoryFunc {
	return func(context es.AggregateContext) es.AggregateRoot {
		return &UserAggregate{
			AggregateContext:       context,
			emailUniquenessService: emailUniquenessService,
		}
	}
}

func (a *UserAggregate) ApplyEvent(event es.EventData) {
	switch e := event.(type) {
	case UserCreatedEvent:
		a.applyUserCreatedEvent(e)
	case UserLoginTokenGeneratedEvent:
		// no-op
	case UserLoginEvent:
		a.applyUserLoginEvent(e)
	default:
		panic("unknown event type")
	}
}

func (a *UserAggregate) HandleCommand(ctx context.Context, command es.Command) ([]es.EventData, error) {
	switch c := command.(type) {
	case CreateUserCommand:
		return a.handleCreateUserCommand(ctx, c)
	case GenerateLoginTokenCommand:
		return a.handleGenerateLoginTokenCommand(ctx, c)
	case LoginCommand:
		return a.handleLoginCommand(ctx, c)
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

	if unique, err := a.emailUniquenessService.IsUnique(ctx, email); err != nil {
		return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
	} else if !unique {
		return nil, fmt.Errorf("email is not unique: %s", command.Email)
	}

	return c.Events(UserCreatedEvent{
		UserID:    userId.String(),
		FirstName: firstName.String(),
		LastName:  lastName.String(),
		Email:     email.String(),
	})
}

func (a *UserAggregate) handleGenerateLoginTokenCommand(ctx context.Context, command GenerateLoginTokenCommand) ([]es.EventData, error) {
	if a.Version() == initialAggregateVersion {
		return nil, fmt.Errorf("user with provided ID does not exists: %s", command.UserID)
	}

	return c.Events(UserLoginTokenGeneratedEvent{
		UserID: a.ID.String(),
		Email:  a.Email.String(),
		Token:  uuid.NewString(),
	})
}

func (a *UserAggregate) handleLoginCommand(ctx context.Context, command LoginCommand) ([]es.EventData, error) {
	if a.Version() == initialAggregateVersion {
		return nil, fmt.Errorf("user with provided ID does not exists: %s", command.UserID)
	}

	loginItem, err := NewLastLogin(time.Now().Unix(), command.UserAgent, command.IP)
	if err != nil {
		return nil, err
	}

	return c.Events(UserLoginEvent{
		UserID:    a.ID.String(),
		Timestamp: loginItem.Time.Unix(),
		UserAgent: loginItem.UserAgent.String(),
		IP:        loginItem.IP.String(),
	})
}

func (a *UserAggregate) applyUserCreatedEvent(event UserCreatedEvent) {
	a.ID, _ = c.NewUserID(event.UserID)
	a.FirstName, _ = NewFirstName(event.FirstName)
	a.LastName, _ = NewLastName(event.LastName)
	a.Email, _ = NewEmail(event.Email)
}

func (a *UserAggregate) applyUserLoginEvent(event UserLoginEvent) {
	a.LastLogin, _ = NewLastLogin(event.Timestamp, event.UserAgent, event.IP)
}
