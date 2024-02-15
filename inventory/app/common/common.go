package common

import (
	"context"

	es "github.com/cybre/home-inventory/pkg/eventsourcing"
)

type CommandBus interface {
	Dispatch(ctx context.Context, command es.Command) error
}
