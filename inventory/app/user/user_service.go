package user

import (
	"context"
	"fmt"

	"github.com/cybre/home-inventory/inventory/app/common"
	"github.com/cybre/home-inventory/inventory/domain/user"
	"github.com/cybre/home-inventory/inventory/shared"
)

type UserIDRepository interface {
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
}

type TokenVerifier interface {
	VerifyToken(ctx context.Context, token string) (string, error)
}

type UserService struct {
	commandBus       common.CommandBus
	userIdRepository UserIDRepository
	tokenVerifier    TokenVerifier
}

func NewUserService(commandBus common.CommandBus, userIdRepository UserIDRepository, tokenVerifier TokenVerifier) *UserService {
	return &UserService{
		commandBus:       commandBus,
		userIdRepository: userIdRepository,
		tokenVerifier:    tokenVerifier,
	}
}

func (s UserService) CreateUser(ctx context.Context, data shared.CreateUserCommandData) error {
	return s.commandBus.Dispatch(ctx, user.CreateUserCommand{
		UserID:    data.UserID,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
	})
}

func (s UserService) GenerateLoginToken(ctx context.Context, data shared.GenerateLoginTokenCommandData) error {
	userId, err := s.userIdRepository.GetUserIDByEmail(ctx, data.Email)
	if err != nil {
		return fmt.Errorf("failed to get user id by email: %w", err)
	}

	return s.commandBus.Dispatch(ctx, user.GenerateLoginTokenCommand{
		UserID: userId,
	})
}

func (s UserService) LoginViaToken(ctx context.Context, data shared.LoginViaTokenCommandData) error {
	userId, err := s.tokenVerifier.VerifyToken(ctx, data.Token)
	if err != nil {
		return fmt.Errorf("failed to get user id by token: %w", err)
	}

	return s.commandBus.Dispatch(ctx, user.LoginCommand{
		UserID:    userId,
		UserAgent: data.UserAgent,
		IP:        data.IP,
	})
}
