package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pkg/errors"

	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq/models"
)

// BetCalculatedConsumer consumes calculated bets from the desired RabbitMQ queue.
type BetCalculatedConsumer struct {
	channel Channel
	config  ConsumerConfig
}

// NewBetCalculatedConsumer creates and returns a new BetCalculatedConsumer.
func NewBetCalculatedConsumer(channel Channel, config ConsumerConfig) (*BetCalculatedConsumer, error) {
	_, err := channel.QueueDeclare(
		config.Queue,
		config.DeclareDurable,
		config.DeclareAutoDelete,
		config.DeclareExclusive,
		config.DeclareNoWait,
		config.DeclareArgs,
	)
	if err != nil {
		return nil, errors.Wrap(err, "bet calculated consumer initialization failed")
	}

	return &BetCalculatedConsumer{
		channel: channel,
		config:  config,
	}, nil
}

// Consume consumes messages until the context is cancelled. An error will be returned if consuming
// is not possible.
func (c *BetCalculatedConsumer) Consume(ctx context.Context) (<-chan models.BetCalculated, error) {
	msgs, err := c.channel.Consume(
		c.config.Queue,
		c.config.ConsumerName,
		c.config.AutoAck,
		c.config.Exclusive,
		c.config.NoLocal,
		c.config.NoWait,
		c.config.Args,
	)
	if err != nil {
		return nil, errors.Wrap(err, "bet calculated consumer failed to consume messages")
	}

	betsCalculated := make(chan models.BetCalculated)
	go func() {
		defer close(betsCalculated)
		for msg := range msgs {
			var betCalculated models.BetCalculated
			err := json.Unmarshal(msg.Body, &betCalculated)
			if err != nil {
				log.Println("Failed to unmarshal bet calculated message", msg.Body)
			}

			// Once context is cancelled, stop consuming messages.
			select {
			case betsCalculated <- betCalculated:
			case <-ctx.Done():
				return
			}
		}
	}()

	return betsCalculated, nil
}
