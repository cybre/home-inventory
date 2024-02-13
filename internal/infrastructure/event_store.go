package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	es "github.com/cybre/home-inventory/pkg/eventsourcing"
	"github.com/gocql/gocql"
)

type CassandraEventStore struct {
	session *gocql.Session
}

func NewCassandraEventStore(session *gocql.Session) *CassandraEventStore {
	return &CassandraEventStore{
		session: session,
	}
}

func (ces *CassandraEventStore) StoreEvents(ctx context.Context, events []es.Event) error {
	batch := ces.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	for _, event := range events {
		eventData, err := json.Marshal(event.Data)
		if err != nil {
			return err
		}

		aggregateID, err := gocql.ParseUUID(string(event.AggregateID))
		if err != nil {
			return err
		}

		batch.Query(
			"INSERT INTO event_store (aggregate_type, aggregate_id, event_type, event_data, timestamp, version) VALUES (?, ?, ?, ?, ?, ?) IF NOT EXISTS",
			event.AggregateType,
			aggregateID,
			event.EventType,
			eventData,
			event.Timestamp,
			event.Version,
		)
	}

	return ces.session.ExecuteBatch(batch)
}

func (ces *CassandraEventStore) GetEvents(aggregateType es.AggregateType, aggregateID es.AggregateID) ([]es.Event, error) {
	aggregateUUID, err := gocql.ParseUUID(string(aggregateID))
	if err != nil {
		return nil, fmt.Errorf("failed to parse aggregate id: %w", err)
	}

	scanner := ces.session.Query(
		"SELECT event_type, event_data, timestamp, version FROM event_store WHERE aggregate_type = ? AND aggregate_id = ?",
		aggregateType,
		aggregateUUID,
	).Iter().Scanner()

	events := []es.Event{}
	for scanner.Next() {
		var (
			eventType string
			eventData []byte
			timestamp int64
			version   uint
		)

		if err := scanner.Scan(&eventType, &eventData, &timestamp, &version); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		eventDataInstance, ok := es.GetEvent(es.EventType(eventType))
		if !ok {
			return nil, es.ErrEventTypeNotFound
		}

		if err := json.Unmarshal(eventData, eventDataInstance); err != nil {
			return nil, fmt.Errorf("failed to decode event data: %w", err)
		}

		eventDataInstanceValue := reflect.ValueOf(eventDataInstance).Elem().Interface().(es.EventData)

		events = append(events, es.Event{
			AggregateType: aggregateType,
			AggregateID:   aggregateID,
			EventType:     es.EventType(eventType),
			Data:          eventDataInstanceValue,
			Timestamp:     timestamp,
			Version:       version,
		})
	}

	return events, nil
}
