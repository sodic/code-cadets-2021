package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq/models"
)

// BetPublisher publishes bets into the desired RabbitMQ queue.
type BetPublisher struct {
	channel Channel
	config  PublisherConfig
}

// NewBetPublisher creates and returns a new BetPublisher.
func NewBetPublisher(channel Channel, config PublisherConfig) (*BetPublisher, error) {
	_, err := channel.QueueDeclare(
		config.Queue,
		config.DeclareDurable,
		config.DeclareAutoDelete,
		config.DeclareExclusive,
		config.DeclareNoWait,
		config.DeclareArgs,
	)
	if err != nil {
		return nil, errors.Wrap(err, "bet publisher initialization failed")
	}

	return &BetPublisher{
		channel: channel,
		config:  config,
	}, nil
}

// Publish publishes messages until the context is cancelled.
func (p *BetPublisher) Publish(ctx context.Context, bets <-chan models.Bet) {
	go func() {
		for bet := range bets {
			select {
			case <-ctx.Done():
				return
			default:
				betJson, err := json.Marshal(&bet)
				if err != nil {
					log.Println("Failed to marshal the following bet:", bet)
					continue
				}

				err = p.channel.Publish(
					p.config.Exchange,
					p.config.Queue,
					p.config.Mandatory,
					p.config.Immediate,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        betJson,
					},
				)
				if err != nil {
					log.Println("Failed to publish the following bet:", betJson)
				}
			}
		}
	}()
}
