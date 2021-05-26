package publisher

import (
	"context"

	rabbitmqmodels "github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/internal/infrastructure/rabbitmq/models"
)

// Publisher offers methods for publishing into output queues.
type Publisher struct {
	betPublisher BetPublisher
}

// New creates and returns a new Publisher.
func New(betPublisher BetPublisher) *Publisher {
	return &Publisher{
		betPublisher: betPublisher,
	}
}

// PublishBets publishes into bets queue.
func (p *Publisher) PublishBets(ctx context.Context, bets <-chan rabbitmqmodels.Bet) {
	p.betPublisher.Publish(ctx, bets)
}
