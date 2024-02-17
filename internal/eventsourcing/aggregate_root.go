package eventsourcing

import "context"

type AggregateID string

func (id AggregateID) Marshal() []byte {
	return []byte(id)
}

func (id AggregateID) String() string {
	return string(id)
}

type AggregateType string

// AggregateRoot represents the interface that all aggregates in the system should implement.
type AggregateRoot interface {
	AggregateID() AggregateID
	AggregateType() AggregateType
	Version() uint
	ApplyEvent(EventData)
	HandleCommand(context.Context, Command) ([]EventData, error)
}

type AggregateRootFactoryFunc func(context AggregateContext) AggregateRoot

type AggregateRootRegistry struct {
	aggregates map[AggregateType]AggregateRootFactoryFunc
}

func NewAggregateRegistry() *AggregateRootRegistry {
	return &AggregateRootRegistry{
		aggregates: make(map[AggregateType]AggregateRootFactoryFunc),
	}
}

func (r *AggregateRootRegistry) RegisterAggregateRoot(aggregateType AggregateType, aggregateFactory AggregateRootFactoryFunc) {
	r.aggregates[aggregateType] = aggregateFactory
}

func (r *AggregateRootRegistry) GetAggregateRoot(aggregateContext AggregateContext) (AggregateRoot, bool) {
	aggregateFactory, ok := r.aggregates[aggregateContext.AggregateType()]
	if ok {
		return aggregateFactory(aggregateContext), true
	}

	return nil, false
}

var aggregateRegistry *AggregateRootRegistry

func init() {
	aggregateRegistry = NewAggregateRegistry()
}

func RegisterAggregateRoot(aggregateType AggregateType, aggregateFactory AggregateRootFactoryFunc) {
	aggregateRegistry.RegisterAggregateRoot(aggregateType, aggregateFactory)
}

func GetAggregateRoot(aggregateContext AggregateContext) (AggregateRoot, bool) {
	return aggregateRegistry.GetAggregateRoot(aggregateContext)
}
