package eventsourcing

type AggregateContext struct {
	aggregateType AggregateType
	aggregateID   AggregateID
	version       uint
}

func NewAggregateContext(aggregateType AggregateType, aggregateID AggregateID, version uint) AggregateContext {
	return AggregateContext{
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
		version:       version,
	}
}

func (a AggregateContext) AggregateID() AggregateID {
	return a.aggregateID
}

func (a AggregateContext) AggregateType() AggregateType {
	return a.aggregateType
}

func (a AggregateContext) Version() uint {
	return a.version
}
