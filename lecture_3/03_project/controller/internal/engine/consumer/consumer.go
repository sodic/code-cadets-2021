package consumer

import (
	"context"

	rabbitmqmodels "github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq/models"
)

// Consumer offers methods for consuming from input queues.
type Consumer struct {
	betReceivedConsumer   BetReceivedConsumer
	betCalculatedConsumer BetCalculatedConsumer
}

// New creates and returns a new Consumer.
func New(betReceivedConsumer BetReceivedConsumer, betCalculatedConsumer BetCalculatedConsumer) *Consumer {
	return &Consumer{
		betReceivedConsumer:   betReceivedConsumer,
		betCalculatedConsumer: betCalculatedConsumer,
	}
}

// ConsumeBetsReceived consumes bets received queue.
func (c *Consumer) ConsumeBetsReceived(ctx context.Context) (<-chan rabbitmqmodels.BetReceived, error) {
	return c.betReceivedConsumer.Consume(ctx)
}

// ConsumeBetsCalculated consumes bets calculated queue.
func (c *Consumer) ConsumeBetsCalculated(ctx context.Context) (<-chan rabbitmqmodels.BetCalculated, error) {
	return c.betCalculatedConsumer.Consume(ctx)
}
