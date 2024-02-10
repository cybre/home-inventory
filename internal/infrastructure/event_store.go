package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/cybre/home-inventory/pkg/domain"
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

func (es *CassandraEventStore) StoreEvents(ctx context.Context, events []domain.Event) error {
	batch := es.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
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

	return es.session.ExecuteBatch(batch)
}

func (es *CassandraEventStore) GetEvents(aggregateType domain.AggregateType, aggregateID domain.AggregateID) ([]domain.Event, error) {
	aggregateUUID, err := gocql.ParseUUID(string(aggregateID))
	if err != nil {
		return nil, fmt.Errorf("failed to parse aggregate id: %w", err)
	}

	scanner := es.session.Query(
		"SELECT event_type, event_data, timestamp, version FROM event_store WHERE aggregate_type = ? AND aggregate_id = ?",
		aggregateType,
		aggregateUUID,
	).Iter().Scanner()

	events := []domain.Event{}
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

		eventDataInstance, ok := domain.GetEvent(domain.EventType(eventType))
		if !ok {
			return nil, domain.ErrEventTypeNotFound
		}

		if err := json.Unmarshal(eventData, eventDataInstance); err != nil {
			return nil, fmt.Errorf("failed to decode event data: %w", err)
		}

		eventDataInstanceValue := reflect.ValueOf(eventDataInstance).Elem().Interface().(domain.EventData)

		events = append(events, domain.Event{
			AggregateType: aggregateType,
			AggregateID:   aggregateID,
			EventType:     domain.EventType(eventType),
			Data:          eventDataInstanceValue,
			Timestamp:     timestamp,
			Version:       version,
		})
	}

	return events, nil
}
