package user

import (
	"context"

	"github.com/cybre/home-inventory/internal/app/common"
	"github.com/cybre/home-inventory/internal/shared"
)

type UserService struct {
	commandBus common.CommandBus
}

func NewUserService(commandBus common.CommandBus) *UserService {
	return &UserService{
		commandBus: commandBus,
	}
}

func (s UserService) CreateUser(ctx context.Context, data shared.CreateUserCommandData) error {
	return s.commandBus.Dispatch(ctx, CreateUserCommand{
		UserID:    data.UserID,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
	})
}
