package app

import (
	"context"

	"github.com/cybre/home-inventory/pkg/domain"
)

type CommandBus interface {
	Dispatch(ctx context.Context, command domain.Command) error
}

type ItemService struct {
	CommandBus CommandBus
}

func NewItemService(commandBus CommandBus) *ItemService {
	return &ItemService{
		CommandBus: commandBus,
	}
}

type AddItemCommandData struct {
	ItemID string `json:"item_id"`
	Name   string `json:"name"`
}

func (s *ItemService) AddItem(ctx context.Context, data AddItemCommandData) error {
	return s.CommandBus.Dispatch(ctx, AddItemCommand{
		ItemID: data.ItemID,
		Name:   data.Name,
	})
}

type UpdateItemCommandData struct {
	ItemID string `json:"item_id"`
	Name   string `json:"name"`
}

func (s *ItemService) UpdateItem(ctx context.Context, data UpdateItemCommandData) error {
	return s.CommandBus.Dispatch(ctx, UpdateItemCommand{
		ItemID: data.ItemID,
		Name:   data.Name,
	})
}
