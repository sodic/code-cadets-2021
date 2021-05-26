package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pkg/errors"

	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq/models"
)

// BetReceivedConsumer consumes received bets from the desired RabbitMQ queue.
type BetReceivedConsumer struct {
	channel Channel
	config  ConsumerConfig
}

// NewBetReceivedConsumer creates and returns a new BetReceivedConsumer.
func NewBetReceivedConsumer(channel Channel, config ConsumerConfig) (*BetReceivedConsumer, error) {
	_, err := channel.QueueDeclare(
		config.Queue,
		config.DeclareDurable,
		config.DeclareAutoDelete,
		config.DeclareExclusive,
		config.DeclareNoWait,
		config.DeclareArgs,
	)
	if err != nil {
		return nil, errors.Wrap(err, "bet received consumer initialization failed")
	}

	return &BetReceivedConsumer{
		channel: channel,
		config:  config,
	}, nil
}

// Consume consumes messages until the context is cancelled. An error will be returned if consuming
// is not possible.
func (c *BetReceivedConsumer) Consume(ctx context.Context) (<-chan models.BetReceived, error) {
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
		return nil, errors.Wrap(err, "bet received consumer failed to consume messages")
	}

	betsReceived := make(chan models.BetReceived)
	go func() {
		defer close(betsReceived)
		for msg := range msgs {
			var betReceived models.BetReceived
			err := json.Unmarshal(msg.Body, &betReceived)
			if err != nil {
				log.Println("Failed to unmarshal bet received message", msg.Body)
			}

			// Once context is cancelled, stop consuming messages.
			select {
			case betsReceived <- betReceived:
			case <-ctx.Done():
				return
			}
		}
	}()

	return betsReceived, nil
}
