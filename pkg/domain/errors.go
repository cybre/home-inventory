package domain

import "errors"

var (
	ErrAggregateTypeNotFound = errors.New("aggregate type not found in registry")
	ErrEventTypeNotFound     = errors.New("event type not found in registry")
)
