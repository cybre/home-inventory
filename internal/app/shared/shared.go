package shared

import (
	"context"

	"github.com/cybre/home-inventory/pkg/domain"
)

const (
	InventoryAggregateType domain.AggregateType = "InventoryAggregate"
)

func Events(events ...domain.EventData) ([]domain.EventData, error) {
	return events, nil
}

type CommandBus interface {
	Dispatch(ctx context.Context, command domain.Command) error
}
