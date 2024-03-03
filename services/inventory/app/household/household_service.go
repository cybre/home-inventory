package household

import (
	"context"
	"fmt"

	"github.com/cybre/home-inventory/internal/utils"
	"github.com/cybre/home-inventory/services/inventory/app/common"
	"github.com/cybre/home-inventory/services/inventory/domain/household"
	"github.com/cybre/home-inventory/services/inventory/shared"
)

type UserHouseholdRepo interface {
	GetUserHouseholds(ctx context.Context, userID string) ([]UserHouseholdModel, error)
	GetUserHousehold(ctx context.Context, userID, householdID string) (UserHouseholdModel, error)
	GetRoom(ctx context.Context, userID, householdID, roomID string) (UserHouseholdRoomModel, error)
}

type HouseholdService struct {
	commandBus common.CommandBus
	repository UserHouseholdRepo
}

func NewHouseholdService(commandBus common.CommandBus, repository UserHouseholdRepo) *HouseholdService {
	return &HouseholdService{
		commandBus: commandBus,
		repository: repository,
	}
}

func (s HouseholdService) CreateHousehold(ctx context.Context, data shared.CreateHouseholdCommandData) error {
	households, err := s.repository.GetUserHouseholds(ctx, data.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user households: %w", err)
	}

	for _, household := range households {
		if household.Name == data.Name {
			return fmt.Errorf("household with name %s already exists", data.Name)
		}
	}

	return s.commandBus.Dispatch(ctx, household.CreateHouseholdCommand{
		HouseholdID: data.HouseholdID,
		UserID:      data.UserID,
		Name:        data.Name,
		Location:    data.Location,
		Description: data.Description,
		Order:       uint(len(households)) + 1,
	})
}

func (s HouseholdService) UpdateHousehold(ctx context.Context, data shared.UpdateHouseholdCommandData) error {
	return s.commandBus.Dispatch(ctx, household.UpdateHouseholdCommand{
		HouseholdID: data.HouseholdID,
		UserID:      data.UserID,
		Name:        data.Name,
		Location:    data.Location,
		Description: data.Description,
	})
}

func (s HouseholdService) DeleteHousehold(ctx context.Context, data shared.DeleteHouseholdCommandData) error {
	return s.commandBus.Dispatch(ctx, household.DeleteHouseholdCommand{
		HouseholdID: data.HouseholdID,
		UserID:      data.UserID,
	})
}

func (s HouseholdService) AddRoom(ctx context.Context, data shared.AddRoomCommandData) error {
	return s.commandBus.Dispatch(ctx, household.AddRoomCommand{
		HouseholdID: data.HouseholdID,
		UserID:      data.UserID,
		RoomID:      data.RoomID,
		Name:        data.Name,
	})
}

func (s HouseholdService) UpdateRoom(ctx context.Context, data shared.UpdateRoomCommandData) error {
	return s.commandBus.Dispatch(ctx, household.UpdateRoomCommand{
		HouseholdID: data.HouseholdID,
		UserID:      data.UserID,
		RoomID:      data.RoomID,
		Name:        data.Name,
	})
}

func (s HouseholdService) DeleteRoom(ctx context.Context, data shared.DeleteRoomCommandData) error {
	return s.commandBus.Dispatch(ctx, household.DeleteRoomCommand{
		HouseholdID: data.HouseholdID,
		UserID:      data.UserID,
		RoomID:      data.RoomID,
	})
}

func (s HouseholdService) AddItem(ctx context.Context, data shared.AddItemCommandData) error {
	return s.commandBus.Dispatch(ctx, household.AddItemCommand{
		HouseholdID: data.HouseholdID,
		RoomID:      data.RoomID,
		ItemID:      data.ItemID,
		Name:        data.Name,
		Barcode:     data.Barcode,
		Quantity:    data.Quantity,
	})
}

func (s HouseholdService) UpdateItem(ctx context.Context, data shared.UpdateItemCommandData) error {
	return s.commandBus.Dispatch(ctx, household.UpdateItemCommand{
		HouseholdID: data.HouseholdID,
		RoomID:      data.RoomID,
		ItemID:      data.ItemID,
		Name:        data.Name,
		Barcode:     data.Barcode,
		Quantity:    data.Quantity,
	})
}

func (s HouseholdService) GetUserHouseholds(ctx context.Context, userID string) ([]shared.UserHousehold, error) {
	households, err := s.repository.GetUserHouseholds(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toSharedUserHouseholds(households), nil
}

func (s HouseholdService) GetUserHousehold(ctx context.Context, userID, householdID string) (shared.UserHousehold, error) {
	household, err := s.repository.GetUserHousehold(ctx, userID, householdID)
	if err != nil {
		return shared.UserHousehold{}, err
	}

	return toSharedUserHousehold(household), nil
}

func (s HouseholdService) GetUserHouseholdRoom(ctx context.Context, userID, householdID, roomID string) (shared.UserHouseholdRoom, error) {
	room, err := s.repository.GetRoom(ctx, userID, householdID, roomID)
	if err != nil {
		return shared.UserHouseholdRoom{}, err
	}

	return toSharedUserHouseholdRoom(0, room), nil
}

func toSharedUserHouseholds(households []UserHouseholdModel) []shared.UserHousehold {
	sharedHouseholds := make([]shared.UserHousehold, len(households))
	for i, household := range households {
		sharedHouseholds[i] = toSharedUserHousehold(household)
	}

	return sharedHouseholds
}

func toSharedUserHousehold(household UserHouseholdModel) shared.UserHousehold {
	return shared.UserHousehold{
		UserID:      household.UserID,
		HouseholdID: household.HouseholdID.String(),
		Name:        household.Name,
		Location:    household.Location,
		Description: household.Description,
		Rooms:       utils.Map(household.Rooms, toSharedUserHouseholdRoom),
		Timestamp:   household.Timestamp,
		Order:       household.Order,
	}
}

func toSharedUserHouseholdRoom(i uint, room UserHouseholdRoomModel) shared.UserHouseholdRoom {
	return shared.UserHouseholdRoom{
		HouseholdID: room.HouseholdID.String(),
		RoomID:      room.RoomID.String(),
		Name:        room.Name,
		Timestamp:   room.Timestamp,
		Order:       room.Order,
	}
}
