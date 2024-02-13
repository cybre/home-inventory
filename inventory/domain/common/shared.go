package common

import es "github.com/cybre/home-inventory/pkg/eventsourcing"

func Events(events ...es.EventData) ([]es.EventData, error) {
	return events, nil
}
