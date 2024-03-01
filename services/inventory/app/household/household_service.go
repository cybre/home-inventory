package household

import (
	"context"

	"github.com/cybre/home-inventory/internal/utils"
	"github.com/cybre/home-inventory/services/inventory/app/common"
	"github.com/cybre/home-inventory/services/inventory/domain/household"
	"github.com/cybre/home-inventory/services/inventory/shared"
)

type UserHouseholdGetter interface {
	GetUserHouseholds(ctx context.Context, userID string) ([]UserHouseholdModel, error)
	GetUserHousehold(ctx context.Context, userID, householdID string) (UserHouseholdModel, error)
}

type HouseholdService struct {
	commandBus common.CommandBus
	repository UserHouseholdGetter
}

func NewHouseholdService(commandBus common.CommandBus, repository UserHouseholdGetter) *HouseholdService {
	return &HouseholdService{
		commandBus: commandBus,
		repository: repository,
	}
}

func (s HouseholdService) CreateHousehold(ctx context.Context, data shared.CreateHouseholdCommandData) error {
	return s.commandBus.Dispatch(ctx, household.CreateHouseholdCommand{
		HouseholdID: data.HouseholdID,
		UserID:      data.UserID,
		Name:        data.Name,
		Location:    data.Location,
		Description: data.Description,
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

func (s HouseholdService) AddRoom(ctx context.Context, data shared.AddRoomCommandData) error {
	return s.commandBus.Dispatch(ctx, household.AddRoomCommand{
		HouseholdID: data.HouseholdID,
		RoomID:      data.RoomID,
		Name:        data.Name,
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
		HouseholdID: household.HouseholdID,
		Name:        household.Name,
		Location:    household.Location,
		Description: household.Description,
		ItemCount:   household.ItemCount,
		Rooms:       utils.Map(household.Rooms, toSharedUserHouseholdRoom),
		Timestamp:   household.Timestamp,
	}
}

func toSharedUserHouseholdRoom(i uint, room UserHouseholdRoomModel) shared.UserHouseholdRoom {
	return shared.UserHouseholdRoom{
		HouseholdID: room.HouseholdID,
		RoomID:      room.RoomID,
		Name:        room.Name,
		ItemCount:   room.ItemCount,
	}
}
