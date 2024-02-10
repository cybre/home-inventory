package items

import (
	"context"

	"github.com/cybre/home-inventory/internal/app/shared"
)

type ItemService struct {
	CommandBus shared.CommandBus
}

func NewItemService(commandBus shared.CommandBus) *ItemService {
	return &ItemService{
		CommandBus: commandBus,
	}
}

type AddItemCommandData struct {
	ItemID string `json:"item_id"`
	Name   string `json:"name"`
}

func (s ItemService) AddItem(ctx context.Context, data AddItemCommandData) error {
	return s.CommandBus.Dispatch(ctx, AddItemCommand{
		ItemID: data.ItemID,
		Name:   data.Name,
	})
}

type UpdateItemCommandData struct {
	ItemID string `json:"item_id"`
	Name   string `json:"name"`
}

func (s ItemService) UpdateItem(ctx context.Context, data UpdateItemCommandData) error {
	return s.CommandBus.Dispatch(ctx, UpdateItemCommand{
		ItemID: data.ItemID,
		Name:   data.Name,
	})
}
