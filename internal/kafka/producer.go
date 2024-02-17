package kafka

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Record struct {
	Key   []byte
	Value []byte
}

type Producer struct {
	client *kgo.Client
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	producerId := strconv.FormatInt(int64(os.Getpid()), 10)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.TransactionalID(producerId),
		kgo.DefaultProduceTopic(topic),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing Kafka producer: %w", err)
	}

	return &Producer{client: cl}, nil
}

func (p Producer) Produce(ctx context.Context, records ...Record) error {
	if err := p.client.BeginTransaction(); err != nil {
		return fmt.Errorf("unable to start transaction: %v", err)
	}

	e := kgo.AbortingFirstErrPromise(p.client)

	for _, record := range records {
		p.client.Produce(ctx, &kgo.Record{
			Key:   record.Key,
			Value: record.Value,
		}, e.Promise())
	}

	commit := kgo.TransactionEndTry(e.Err() == nil)

	switch err := p.client.EndTransaction(ctx, commit); err {
	case nil:
	case kerr.OperationNotAttempted:
		if err := p.client.EndTransaction(ctx, kgo.TryAbort); err != nil {
			return fmt.Errorf("unable to abort transaction: %v", err)
		}
	default:
		return fmt.Errorf("unable to end transaction: %v", err)
	}

	return nil
}

func (p Producer) Close() {
	p.client.Close()
}
