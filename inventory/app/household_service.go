package app

import (
	"context"

	"github.com/cybre/home-inventory/inventory/domain/household"
	"github.com/cybre/home-inventory/inventory/shared"
)

type HouseholdService struct {
	CommandBus CommandBus
}

func NewHouseholdService(commandBus CommandBus) *HouseholdService {
	return &HouseholdService{
		CommandBus: commandBus,
	}
}

func (s HouseholdService) CreateHousehold(ctx context.Context, data shared.CreateHouseholdCommandData) error {
	return s.CommandBus.Dispatch(ctx, household.CreateHouseholdCommand{
		HouseholdID: data.HouseholdID,
		UserID:      data.UserID,
		Name:        data.Name,
	})
}

func (s HouseholdService) AddRoom(ctx context.Context, data shared.AddRoomCommandData) error {
	return s.CommandBus.Dispatch(ctx, household.AddRoomCommand{
		HouseholdID: data.HouseholdID,
		RoomID:      data.RoomID,
		Name:        data.Name,
	})
}

func (s HouseholdService) AddItem(ctx context.Context, data shared.AddItemCommandData) error {
	return s.CommandBus.Dispatch(ctx, household.AddItemCommand{
		ItemID:   data.ItemID,
		Name:     data.Name,
		Barcode:  data.Barcode,
		Quantity: data.Quantity,
	})
}

func (s HouseholdService) UpdateItem(ctx context.Context, data shared.UpdateItemCommandData) error {
	return s.CommandBus.Dispatch(ctx, household.UpdateItemCommand{
		ItemID:   data.ItemID,
		Name:     data.Name,
		Barcode:  data.Barcode,
		Quantity: data.Quantity,
	})
}
