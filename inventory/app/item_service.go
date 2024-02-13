package app

import (
	"context"

	"github.com/cybre/home-inventory/inventory/domain/user"
	"github.com/cybre/home-inventory/inventory/shared"
)

type UserService struct {
	commandBus CommandBus
}

func NewUserService(commandBus CommandBus) *UserService {
	return &UserService{
		commandBus: commandBus,
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
