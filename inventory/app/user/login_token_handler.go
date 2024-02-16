package user

import (
	"context"
	"fmt"

	es "github.com/cybre/home-inventory/pkg/eventsourcing"
	"github.com/cybre/home-inventory/pkg/logging"

	domainuser "github.com/cybre/home-inventory/inventory/domain/user"
)

type LoginTokenInserter interface {
	InsertToken(ctx context.Context, userId, token string) error
}

type LoginTokenHandler struct {
	repository LoginTokenInserter
}

func NewLoginTokenHandler(repository LoginTokenInserter) *LoginTokenHandler {
	return &LoginTokenHandler{
		repository: repository,
	}
}

func (h LoginTokenHandler) HandleEvent(ctx context.Context, event es.EventData) error {
	switch e := event.(type) {
	case domainuser.UserLoginTokenGeneratedEvent:
		return h.handleUserLoginTokenGeneratedEvent(ctx, e)
	default:
		return es.ErrUnknownEvent
	}
}

func (h LoginTokenHandler) Events() []es.EventType {
	return []es.EventType{
		domainuser.EventTypeUserLoginTokenGenerated,
	}
}

func (h LoginTokenHandler) Name() string {
	return "user.LoginTokenHandler"
}

func (h LoginTokenHandler) handleUserLoginTokenGeneratedEvent(ctx context.Context, e domainuser.UserLoginTokenGeneratedEvent) error {
	if err := h.repository.InsertToken(ctx, e.UserID, e.Token); err != nil {
		return fmt.Errorf("failed to insert one time login token for user: %s: %w", e.UserID, err)
	}

	logging.FromContext(ctx).Info("one time login token generated", "userId", e.UserID, "email", e.Email, "token", e.Token)

	return nil
}
