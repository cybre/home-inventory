package user

import (
	"context"
	"fmt"

	es "github.com/cybre/home-inventory/pkg/eventsourcing"
	"github.com/cybre/home-inventory/pkg/logging"

	domainuser "github.com/cybre/home-inventory/inventory/domain/user"
)

type OneTimeLoginTokenInsertRepository interface {
	InsertToken(ctx context.Context, userId, token string) error
}

type OneTimeLoginTokenHandler struct {
	repository OneTimeLoginTokenInsertRepository
}

func NewOneTimeLoginTokenHandler(repository OneTimeLoginTokenInsertRepository) *OneTimeLoginTokenHandler {
	return &OneTimeLoginTokenHandler{
		repository: repository,
	}
}

func (h OneTimeLoginTokenHandler) HandleEvent(ctx context.Context, event es.EventData) error {
	switch e := event.(type) {
	case domainuser.UserOneTimeTokenGeneratedEvent:
		return h.handleOneTimeLoginTokenGeneratedEvent(ctx, e)
	default:
		return es.ErrUnknownEvent
	}
}

func (h OneTimeLoginTokenHandler) Events() []es.EventType {
	return []es.EventType{
		domainuser.EventTypeUserOneTimeTokenGenerated,
	}
}

func (h OneTimeLoginTokenHandler) Name() string {
	return "user.OneTimeLoginTokenHandler"
}

func (h OneTimeLoginTokenHandler) handleOneTimeLoginTokenGeneratedEvent(ctx context.Context, e domainuser.UserOneTimeTokenGeneratedEvent) error {
	if err := h.repository.InsertToken(ctx, e.UserID, e.Token); err != nil {
		return fmt.Errorf("failed to insert one time login token for user: %s: %w", e.UserID, err)
	}

	logging.FromContext(ctx).Info("one time login token generated", "userId", e.UserID, "email", e.Email, "token", e.Token)

	return nil
}
