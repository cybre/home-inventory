package shared

import "github.com/cybre/home-inventory/pkg/domain"

func Events(events ...domain.EventData) ([]domain.EventData, error) {
	return events, nil
}
