package domain

import "errors"

var (
	ErrAggregateTypeNotFound = errors.New("aggregate type not found in registry")
	ErrEventTypeNotFound     = errors.New("event type not found in registry")

	ErrUnknownCommand = errors.New("aggregate does not know how to handle command")
)
