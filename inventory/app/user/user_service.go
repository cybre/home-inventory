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

type UserService struct {
	commandBus       common.CommandBus
	userIdRepository UserIDRepository
}

func NewUserService(commandBus common.CommandBus, userIdRepository UserIDRepository) *UserService {
	return &UserService{
		commandBus:       commandBus,
		userIdRepository: userIdRepository,
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

func (s UserService) GenerateOneTimeToken(ctx context.Context, data shared.GenerateOneTimeTokenCommandData) error {
	userId, err := s.userIdRepository.GetUserIDByEmail(ctx, data.Email)
	if err != nil {
		return fmt.Errorf("failed to get user id by email: %w", err)
	}

	return s.commandBus.Dispatch(ctx, user.GenerateOneTimeTokenCommand{
		UserID: userId,
	})
}
