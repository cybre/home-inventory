package common

import es "github.com/cybre/home-inventory/internal/eventsourcing"

func Events(events ...es.EventData) ([]es.EventData, error) {
	return events, nil
}
