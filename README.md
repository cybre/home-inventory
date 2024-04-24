# home-inventory

A work-in-progress PoC of a home inventory system built with CQRS and event sourcing utlizing domain-driven design.
Kafka is used for event messaging and Cassandra for the event store and projections.  
Client-side caching is handled by Redis.  

## Structure
`internal/eventsourcing` contains the basic building blocks for event sourcing.  
`internal/kafka` contains abstractions for producing and consuming events from kafka topics.  
`internal/infrastructure` contains implementations of a cassandra event store and a kafka event messaging queue.  
`internal/requestbuilder` contains a HTTP request builder (using builder pattern) that supports client-side caching with automatic and manual cache invalidation options