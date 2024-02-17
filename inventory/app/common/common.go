package common

import (
	"context"

	es "github.com/cybre/home-inventory/internal/eventsourcing"
)

type CommandBus interface {
	Dispatch(ctx context.Context, command es.Command) error
}
