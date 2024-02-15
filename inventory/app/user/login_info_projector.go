package user

import (
	"context"
	"fmt"

	"github.com/cybre/home-inventory/inventory/domain/user"
	es "github.com/cybre/home-inventory/pkg/eventsourcing"
)

type LoginInfoInsertRepository interface {
	Insert(ctx context.Context, email, userID string) error
}

type LoginInfoProjector struct {
	repository LoginInfoInsertRepository
}

func NewLoginInfoProjector(repository LoginInfoInsertRepository) *LoginInfoProjector {
	return &LoginInfoProjector{
		repository: repository,
	}
}

func (p LoginInfoProjector) HandleEvent(ctx context.Context, event es.EventData) error {
	switch e := event.(type) {
	case user.UserCreatedEvent:
		return p.handleUserCreatedEvent(ctx, e)
	default:
		return es.ErrUnknownEvent
	}
}

func (p LoginInfoProjector) Events() []es.EventType {
	return []es.EventType{
		user.EventTypeUserCreated,
	}
}

func (p LoginInfoProjector) Name() string {
	return "inventory.LoginInfoProjector"
}

func (p LoginInfoProjector) handleUserCreatedEvent(ctx context.Context, e user.UserCreatedEvent) error {
	if err := p.repository.Insert(ctx, e.Email, e.UserID); err != nil {
		return fmt.Errorf("failed to insert login info for user: %s: %w", e.UserID, err)
	}

	return nil
}
