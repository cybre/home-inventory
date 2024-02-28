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

func toSharedUserHouseholds(households []UserHouseholdModel) []shared.UserHousehold {
	sharedHouseholds := make([]shared.UserHousehold, len(households))
	for i, household := range households {
		sharedHouseholds[i] = shared.UserHousehold{
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

	return sharedHouseholds
}

func toSharedUserHouseholdRoom(i uint, room UserHouseholdRoomModel) shared.UserHouseholdRoom {
	return shared.UserHouseholdRoom{
		RoomID:    room.RoomID,
		Name:      room.Name,
		ItemCount: room.ItemCount,
	}
}
