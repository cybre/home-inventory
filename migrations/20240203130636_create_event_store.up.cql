CREATE TABLE event_store (
    aggregate_type TEXT,
    aggregate_id UUID,
    event_type TEXT,
    event_data TEXT,
    timestamp TIMESTAMP,
    version INT,
    PRIMARY KEY ((aggregate_type, aggregate_id), version)
);