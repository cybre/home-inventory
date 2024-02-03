package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
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
		eventData, err := json.Marshal(event.EventData())
		if err != nil {
			return err
		}

		aggregateID, err := gocql.ParseUUID(string(event.AggregateID()))
		if err != nil {
			return err
		}

		batch.Query(
			"INSERT INTO event_store (aggregate_type, aggregate_id, event_type, event_data, timestamp, version) VALUES (?, ?, ?, ?, ?, ?) IF NOT EXISTS",
			event.AggregateType(),
			aggregateID,
			event.EventType(),
			eventData,
			event.Timestamp(),
			event.Version(),
		)
	}

	return es.session.ExecuteBatch(batch)
}

func (es *CassandraEventStore) GetEvents(aggregateType domain.AggregateType, aggregateID domain.AggregateID) ([]domain.Event, error) {
	aggregateUUID, err := gocql.ParseUUID(string(aggregateID))
	if err != nil {
		return nil, err
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
			return nil, err
		}

		eventDataInstance, ok := domain.GetEvent(domain.EventType(eventType))
		if !ok {
			return nil, domain.ErrEventTypeNotFound
		}

		if err := json.NewDecoder(bytes.NewReader(eventData)).Decode(eventDataInstance); err != nil {
			return nil, err
		}

		eventDataInstanceValue := reflect.ValueOf(eventDataInstance).Elem().Interface().(domain.EventData)

		events = append(events, domain.NewEvent(aggregateType, aggregateID, domain.EventType(eventType), eventDataInstanceValue, timestamp, version))
	}

	return events, nil
}
